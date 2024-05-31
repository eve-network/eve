package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func cosmosnft(contract string, percent int64) ([]banktypes.Balance, []config.Reward, int, error) {
	tokenIds, err := fetchTokenIdsWithRetry(contract)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch token ids: %w", err)
	}

	// Create channels to receive balance info and rewards
	balanceCh := make(chan banktypes.Balance)
	rewardCh := make(chan config.Reward)
	doneCh := make(chan struct{})

	// Use a buffered channel as a semaphore to limit concurrency
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent goroutines

	wg := &sync.WaitGroup{}
	wg.Add(len(tokenIds))

	for _, token := range tokenIds {
		semaphore <- struct{}{} // Acquire semaphore
		go func(token string) {
			defer func() { <-semaphore }() // Release semaphore when done
			defer wg.Done()

			nftHolders, err := fetchTokenInfoWithRetry(token, contract)
			if err != nil {
				fmt.Printf("Error fetching token info for %s: %v\n", token, err)
				return
			}

			if nftHolders.Address == "" {
				fmt.Printf("Token id: %s is not NFT\n", token)
				return
			}

			eveBech32Address, err := utils.ConvertBech32Address(nftHolders.Address)
			if err != nil {
				fmt.Printf("Error converting Bech32Address for %s: %v\n", nftHolders.Address, err)
				return
			}

			allEveAirdrop := sdkmath.LegacyMustNewDecFromStr(EveAirdrop)
			eveAirdrop := (allEveAirdrop.MulInt64(percent)).QuoInt64(100).QuoInt(sdkmath.NewInt(int64(len(tokenIds))))
			rewardCh <- config.Reward{
				Address:         nftHolders.Address,
				EveAddress:      eveBech32Address,
				EveAirdropToken: eveAirdrop,
				ChainID:         config.GetBadKidsConfig().ChainID,
			}

			balanceCh <- banktypes.Balance{
				Address: eveBech32Address,
				Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
			}
		}(token)
	}

	// Close channels once all goroutines are done
	go func() {
		wg.Wait()
		close(balanceCh)
		close(rewardCh)
		close(doneCh)
	}()

	var balanceInfo []banktypes.Balance
	var rewardInfo []config.Reward

	for {
		select {
		case balance, ok := <-balanceCh:
			if !ok {
				balanceCh = nil
			} else {
				balanceInfo = append(balanceInfo, balance)
			}
		case reward, ok := <-rewardCh:
			if !ok {
				rewardCh = nil
			} else {
				rewardInfo = append(rewardInfo, reward)
			}
		case <-doneCh:
			return balanceInfo, rewardInfo, len(balanceInfo), nil
		}
	}
}

func fetchTokenInfoWithRetry(token, contract string) (config.NftHolder, error) {
	var data config.NftHolder
	var err error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		data, err = fetchTokenInfo(token, contract)
		if err == nil {
			return data, nil
		}

		fmt.Printf("Error fetch token info (attempt %d/%d): %v\n", attempt, MaxRetries, err)

		if attempt < MaxRetries {
			// Calculate backoff duration using exponential backoff strategy
			backoffDuration := time.Duration(BackOff.Seconds() * math.Pow(2, float64(attempt)))
			fmt.Printf("Retrying after %s...\n", backoffDuration)
			time.Sleep(backoffDuration)
		}
	}
	return config.NftHolder{}, fmt.Errorf("failed to fetch token info after %d attempts", MaxRetries)
}

func fetchTokenInfo(token, contract string) (config.NftHolder, error) {
	queryString := fmt.Sprintf(`{"all_nft_info":{"token_id":%s}}`, token)
	encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
	apiURL := config.GetStargazeConfig().API + "/cosmwasm/wasm/v1/contract/" + contract + "/smart/" + encodedQuery
	response, err := utils.MakeGetRequest(apiURL)
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

		fmt.Printf("Error fetch token ids (attempt %d/%d): %v\n", attempt, MaxRetries, err)

		if attempt < MaxRetries {
			// Calculate backoff duration using exponential backoff strategy
			backoffDuration := time.Duration(BackOff.Seconds() * math.Pow(2, float64(attempt)))
			fmt.Printf("Retrying after %s...\n", backoffDuration)
			time.Sleep(backoffDuration)
		}
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
		response, err := utils.MakeGetRequest(apiURL)
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
