package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/eve-network/eve/airdrop/config"
)

func cosmosnft(contract string) ([]banktypes.Balance, []config.Reward) {
	tokenIds := fetchTokenIds(contract)
	allEveAirdrop := math.LegacyMustNewDecFromStr(EVE_AIRDROP)
	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}
	testAmount, _ := math.LegacyNewDecFromStr("0")
	eveAirdrop := (allEveAirdrop.MulInt64(int64(config.GetAkashConfig().Percent))).QuoInt64(100).QuoInt(math.NewInt(int64(len(tokenIds))))
	fmt.Println("balance ", eveAirdrop)
	for index, token := range tokenIds {
		nftHolders := fetchTokenInfo(token, contract)
		fmt.Println(index)
		eveBech32Address := convertBech32Address(nftHolders.Address)
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         nftHolders.Address,
			EveAddress:      eveBech32Address,
			EveAirdropToken: eveAirdrop,
			ChainId:         config.GetAkashConfig().ChainID,
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

func fetchTokenInfo(token string, contract string) config.NftHolder {
	queryString := fmt.Sprintf(`{"all_nft_info":{"token_id":%s}}`, token)
	encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
	apiUrl := config.GetStargazeConfig().API + "/cosmwasm/wasm/v1/contract/" + contract + "/smart/" + encodedQuery
	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		panic("")
	}
	defer response.Body.Close()

	var data config.TokenInfoResponse
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		panic("")
	}
	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}
	return config.NftHolder{
		Address: data.Data.Access.Owner,
		TokenId: token,
	}
}

func fetchTokenIds(contract string) []string {
	// Make a GET request to the API
	paginationKey := "0"
	index := 0
	tokenIds := []string{}
	for {
		index += 1
		queryString := fmt.Sprintf(`{"all_tokens":{"limit":1000,"start_after":"%s"}}`, paginationKey)
		encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
		apiUrl := config.GetStargazeConfig().API + "/cosmwasm/wasm/v1/contract/" + contract + "/smart/" + encodedQuery
		response, err := http.Get(apiUrl)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			panic("")
		}
		defer response.Body.Close()

		var data config.TokenIdsResponse
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			panic("")
		}
		// Unmarshal the JSON byte slice into the defined struct
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			panic("")
		}
		tokenIds = append(tokenIds, data.Data.Token...)
		if len(data.Data.Token) == 0 {
			break
		} else {
			paginationKey = data.Data.Token[len(data.Data.Token)-1]
			fmt.Println("pagination key:", paginationKey)
			if len(paginationKey) == 0 {
				break
			}
		}
	}

	fmt.Println(len(tokenIds))
	return tokenIds
}
