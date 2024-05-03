package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eve-network/eve/airdrop/config"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func cosmosnft(contract string, percent int64) ([]banktypes.Balance, []config.Reward, int, error) {
	tokenIds, err := fetchTokenIdsWithRetry(contract)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch token ids: %w", err)
	}
	allEveAirdrop := math.LegacyMustNewDecFromStr(EveAirdrop)
	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}
	testAmount, _ := math.LegacyNewDecFromStr("0")
	eveAirdrop := (allEveAirdrop.MulInt64(percent)).QuoInt64(100).QuoInt(math.NewInt(int64(len(tokenIds))))
	fmt.Println("balance ", eveAirdrop)
	for index, token := range tokenIds {
		nftHolders, err := fetchTokenInfoWithRetry(token, contract)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to fetch token info: %w", err)
		}
		fmt.Println(index)
		eveBech32Address := convertBech32Address(nftHolders.Address)
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         nftHolders.Address,
			EveAddress:      eveBech32Address,
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetBadKidsConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println(testAmount)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}

func fetchTokenInfoWithRetry(token, contract string) (config.NftHolder, error) {
	var data config.NftHolder
	var err error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		data, err = fetchTokenInfo(token, contract)
		if err == nil {
			return data, nil
		}
		fmt.Printf("error fetch token info (attempt %d/%d): %v\n", attempt, MaxRetries, err)
		time.Sleep(time.Duration(time.Duration(attempt * Backoff).Milliseconds()))
	}
	return config.NftHolder{}, fmt.Errorf("failed to fetch token info after %d attempts", MaxRetries)
}

func fetchTokenInfo(token, contract string) (config.NftHolder, error) {
	queryString := fmt.Sprintf(`{"all_nft_info":{"token_id":%s}}`, token)
	encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
	apiURL := config.GetStargazeConfig().API + "/cosmwasm/wasm/v1/contract/" + contract + "/smart/" + encodedQuery
	response, err := http.Get(apiURL)
	if err != nil {
		return config.NftHolder{}, fmt.Errorf("error making GET request to fetch token info: %w", err)
	}
	defer response.Body.Close()

	var data config.TokenInfoResponse
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return config.NftHolder{}, fmt.Errorf("error reading response body when fetch token info: %w", err)
	}
	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return config.NftHolder{}, fmt.Errorf("error unmarshalling JSON when fetch token info: %w", err)
	}
	return config.NftHolder{
		Address: data.Data.Access.Owner,
		TokenID: token,
	}, nil
}

func fetchTokenIdsWithRetry(contract string) ([]string, error) {
	var tokenIds []string
	var err error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		tokenIds, err = fetchTokenIds(contract)
		if err == nil {
			return tokenIds, nil
		}
		fmt.Printf("error fetch token ids (attempt %d/%d): %v\n", attempt, MaxRetries, err)
		time.Sleep(time.Duration(time.Duration(attempt * Backoff).Milliseconds()))
	}
	return nil, fmt.Errorf("failed to fetch token ids after %d attempts", MaxRetries)
}

func fetchTokenIds(contract string) ([]string, error) {
	// Make a GET request to the API
	paginationKey := "0"
	index := 0
	tokenIds := []string{}
	for {
		index++
		queryString := fmt.Sprintf(`{"all_tokens":{"limit":1000,"start_after":"%s"}}`, paginationKey)
		encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
		apiURL := config.GetStargazeConfig().API + "/cosmwasm/wasm/v1/contract/" + contract + "/smart/" + encodedQuery
		response, err := http.Get(apiURL)
		if err != nil {
			return nil, fmt.Errorf("error making GET request to fetch token ids: %w", err)
		}
		defer response.Body.Close()

		var data config.TokenIdsResponse
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body when fetch token ids: %w", err)
		}
		// Unmarshal the JSON byte slice into the defined struct
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			return nil, fmt.Errorf("error error unmarshalling JSON when fetch token ids: %w", err)
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
	return tokenIds, nil
}
