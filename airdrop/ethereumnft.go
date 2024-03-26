package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/eve-network/eve/airdrop/config"
	"github.com/joho/godotenv"
)

const MILADY = "0x5af0d9827e0c53e4799bb226655a1de152a425a5"

func ethereumnft() ([]banktypes.Balance, []config.Reward) {
	nftOwners := fetchNftOwners()
	allEveAirdrop := math.LegacyMustNewDecFromStr(EVE_AIRDROP)
	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}
	testAmount, _ := math.LegacyNewDecFromStr("0")
	eveAirdrop := (allEveAirdrop.MulInt64(int64(config.GetMiladyConfig().Percent))).QuoInt64(100).QuoInt(math.NewInt(int64(len(nftOwners))))
	fmt.Println("balance ", eveAirdrop)
	for index, owner := range nftOwners {
		fmt.Println(index)
		eveBech32Address := convertEvmAddress(owner.OwnerOf)
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         owner.OwnerOf,
			EveAddress:      eveBech32Address,
			EveAirdropToken: eveAirdrop,
			ChainId:         config.GetMiladyConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println(testAmount)
	return balanceInfo, rewardInfo
}

func constructMoralisUrl(cursor string) string {
	return "https://deep-index.moralis.io/api/v2.2/nft/" + MILADY + "/owners?chain=eth&format=decimal&limit=100&cursor=" + cursor
}

func fetchNftOwners() []config.EthResult {
	godotenv.Load()
	API_KEY := os.Getenv("API_KEY")
	pageCount := 0
	cursor := ""
	nftOwners := []config.EthResult{}
	for {
		pageCount += 1
		fmt.Println("Page ", pageCount)
		url := constructMoralisUrl(cursor)
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("Accept", "application/json")
		req.Header.Add("X-API-Key", API_KEY)

		res, _ := http.DefaultClient.Do(req)

		body, _ := ioutil.ReadAll(res.Body)
		var data config.NftEthResponse

		// Unmarshal the JSON byte slice into the defined struct
		err := json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			panic("")
		}
		defer res.Body.Close()

		nftOwners = append(nftOwners, data.Result...)
		if data.Cursor == "" {
			break
		} else {
			cursor = data.Cursor
		}
	}
	return nftOwners
}

func convertEvmAddress(evmAddress string) string {
	addr := common.HexToAddress(evmAddress)
	accCodec := addresscodec.NewBech32Codec("eve")
	eveAddress, err := StringFromEthAddress(accCodec, addr)
	if err != nil {
		fmt.Println("err ", err)
	}
	return eveAddress
}

// EthAddressFromString converts a Cosmos SDK address string to an Ethereum `Address`.
func EthAddressFromString(codec address.Codec, addr string) (common.Address, error) {
	bz, err := codec.StringToBytes(addr)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(bz), nil
}

// MustEthAddressFromString converts a Cosmos SDK address string to an Ethereum `Address`. It
// panics if the conversion fails.
func MustEthAddressFromString(codec address.Codec, addr string) common.Address {
	address, err := EthAddressFromString(codec, addr)
	if err != nil {
		panic(err)
	}
	return address
}

// StringFromEthAddress converts an Ethereum `Address` to a Cosmos SDK address string.
func StringFromEthAddress(codec address.Codec, ethAddress common.Address) (string, error) {
	return codec.BytesToString(ethAddress.Bytes())
}

// MustStringFromEthAddress converts an Ethereum `Address` to a Cosmos SDK address string. It
// panics if the conversion fails.
func MustStringFromEthAddress(codec address.Codec, ethAddress common.Address) string {
	addr, err := StringFromEthAddress(codec, ethAddress)
	if err != nil {
		panic(err)
	}
	return addr
}
