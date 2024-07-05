package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"
	"github.com/joho/godotenv"

	"cosmossdk.io/math"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const MILADY = "0x5af0d9827e0c53e4799bb226655a1de152a425a5"

func Ethereumnft() ([]banktypes.Balance, []config.Reward, int, error) {
	nftOwners, err := fetchNftOwners()
	if err != nil {
		log.Printf("Failed to fetch nft owners: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch nft owners: %w", err)
	}

	allEveAirdrop, err := math.LegacyNewDecFromStr(config.EveAirdrop)
	if err != nil {
		log.Printf("Failed to convert EveAirdrop string to dec: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to convert EveAirdrop string to dec: %w", err)
	}

	if len(nftOwners) == 0 {
		return nil, nil, 0, nil
	}

	eveAirdrop := (allEveAirdrop.MulInt64(int64(config.GetMiladyConfig().Percent))).QuoInt64(100).QuoInt(math.NewInt(int64(len(nftOwners))))
	rewardInfo := make([]config.Reward, len(nftOwners))
	balanceInfo := make([]banktypes.Balance, len(nftOwners))
	totalAmount := math.LegacyMustNewDecFromStr("0")

	for index, owner := range nftOwners {
		eveBech32Address, err := convertEvmAddress(owner.OwnerOf)
		if err != nil {
			log.Printf("Failed to convert evm address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert evm address: %w", err)
		}

		reward := config.Reward{
			Address:         owner.OwnerOf,
			EveAddress:      eveBech32Address,
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetMiladyConfig().ChainID,
		}
		rewardInfo[index] = reward

		balance := banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		}
		balanceInfo[index] = balance

		totalAmount = totalAmount.Add(eveAirdrop)
	}

	log.Println("Total airdrop amount:", totalAmount)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}

func constructMoralisURL(cursor string) string {
	return fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/nft/%s/owners?chain=eth&format=decimal&limit=100&cursor=%s", MILADY, cursor)
}

func fetchNftOwners() ([]config.EthResult, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}

	apiKey := os.Getenv("API_KEY")
	cursor := ""
	var nftOwners []config.EthResult

	for pageCount := 1; ; pageCount++ {
		log.Printf("Fetching page %d", pageCount)

		url := constructMoralisURL(cursor)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("X-API-Key", apiKey)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var data config.NftEthResponse
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}

		nftOwners = append(nftOwners, data.Result...)
		if data.Cursor == "" {
			break
		}
		cursor = data.Cursor
	}

	return nftOwners, nil
}

func convertEvmAddress(evmAddress string) (string, error) {
	addr := common.HexToAddress(evmAddress)
	accCodec := addresscodec.NewBech32Codec("eve")
	return utils.StringFromEthAddress(accCodec, addr)
}
