package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc/metadata"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// got to export genesis state from neutron and bostrom chain

const (
	EveAirdrop   = "1000000000" // 1,000,000,000
	LimitPerPage = 100000000
	Badkids      = "stars19jq6mj84cnt9p7sagjxqf8hxtczwc8wlpuwe4sh62w45aheseues57n420"
	Cryptonium   = "stars1g2ptrqnky5pu70r3g584zpk76cwqplyc63e8apwayau6l3jr8c0sp9q45u"
	APICoingecko = "https://api.coingecko.com/api/v3/simple/price?ids="
	MaxRetries   = 5
	Backoff      = 200 * time.Millisecond
)

// Define a function type that returns balance info, reward info and length
type balanceFunction func() ([]banktypes.Balance, []config.Reward, int, error)

func main() {
	// Define balance functions with their associated names
	balanceFunctions := map[string]balanceFunction{
		"akash":      akash,
		"bostrom":    bostrom,
		"celestia":   celestia,
		"composable": composable,
		"cosmos":     cosmos,
		"neutron":    neutron,
		"sentinel":   sentinel,
		"stargaze":   stargaze,
		"terra":      terra,
		"terrac":     terrac,
		"badkids": func() ([]banktypes.Balance, []config.Reward, int, error) {
			return cosmosnft(Badkids, int64(config.GetBadKidsConfig().Percent))
		},
		"cryptonium": func() ([]banktypes.Balance, []config.Reward, int, error) {
			return cosmosnft(Cryptonium, int64(config.GetCryptoniumConfig().Percent))
		},
		// need set coin type on Eve
		"milady": ethereumnft,
	}

	lenBalanceFunctions := len(balanceFunctions)
	wg := &sync.WaitGroup{}
	wg.Add(lenBalanceFunctions)

	// Channel to collect balance info from goroutines
	balanceInfoCh := make(chan []banktypes.Balance, lenBalanceFunctions)

	// Channel to collect length of balance info from goroutines
	lengthBalanceInfoCh := make(chan int, lenBalanceFunctions)

	// Iterate over the balanceFunctions map and run each function in a goroutine
	for name, fn := range balanceFunctions {
		go func(name string, fn balanceFunction) {
			defer wg.Done()

			fmt.Println("fetching balance info: ", name)
			//TODO: need handle error
			info, _, len, _ := fn()    // Call the function
			balanceInfoCh <- info      // Send balance info to channel
			lengthBalanceInfoCh <- len // Send length of balance info to channel
		}(name, fn)
	}

	go func() {
		// Wait for all goroutines to complete
		wg.Wait()
		// Close channels
		close(balanceInfoCh)
		close(lengthBalanceInfoCh)
	}()

	total := 0
	balanceAkashInfo := []banktypes.Balance{}

	// Collect results from channels
	for lenCh := range lengthBalanceInfoCh {
		total += lenCh
	}

	for infoCh := range balanceInfoCh {
		balanceAkashInfo = append(balanceAkashInfo, infoCh...)
	}

	fmt.Println("total: ", total)
	fmt.Println(len(balanceAkashInfo))

	airdropMap := make(map[string]int)
	for _, info := range balanceAkashInfo {
		amount := airdropMap[info.Address]
		airdropMap[info.Address] = amount + int(info.Coins.AmountOf("eve").Int64())
	}

	balanceInfo := []banktypes.Balance{}
	checkBalance := 0
	for address, amount := range airdropMap {
		checkBalance += amount
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", math.NewInt(int64(amount)))),
		})
	}

	fmt.Println("Check balance: ", checkBalance)

	// // Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardComposableInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	_ = os.WriteFile("balance.json", fileBalance, 0o600)
}

func findValidatorInfo(validators []stakingtypes.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

func getLatestHeightWithRetry(rpcURL string) (string, error) {
	var latestBlockHeight string
	var err error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		latestBlockHeight, err = getLatestHeight(rpcURL)
		if err == nil {
			return latestBlockHeight, nil
		}

		fmt.Printf("error get latest height (attempt %d/%d): %v\n", attempt, MaxRetries, err)

		if attempt < MaxRetries {
			// Calculate backoff duration using exponential backoff strategy
			backoffDuration := time.Duration(Backoff.Seconds() * float64(attempt))
			fmt.Printf("retrying after %s...\n", backoffDuration)
			time.Sleep(backoffDuration)
		}
	}

	return "", fmt.Errorf("failed to get latest height after %d attempts", MaxRetries)
}

func getLatestHeight(apiURL string) (string, error) {
	// Make a GET request to the API
	response, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the response body into a NodeResponse struct
	var data config.NodeResponse
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Extract the latest block height from the response
	latestBlockHeight := data.Result.SyncInfo.LatestBlockHeight
	fmt.Println("Block height:", latestBlockHeight)

	return latestBlockHeight, nil
}

func convertBech32Address(otherChainAddress string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(otherChainAddress)
	if err != nil {
		return "", fmt.Errorf("error decoding address: %w", err)
	}
	newBech32DelAddr, err := bech32.ConvertAndEncode("eve", bz)
	if err != nil {
		return "", fmt.Errorf("error converting address: %w", err)
	}
	return newBech32DelAddr, nil
}

func fetchValidatorsWithRetry(rpcURL string) (config.ValidatorResponse, error) {
	var data config.ValidatorResponse
	var err error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		data, err = fetchValidators(rpcURL)
		if err == nil {
			return data, nil
		}
		fmt.Printf("Error fetching validator info (attempt %d/%d): %v\n", attempt, MaxRetries, err)
		time.Sleep(time.Duration(Backoff.Seconds() * float64(attempt)))
	}
	return config.ValidatorResponse{}, fmt.Errorf("failed to fetch validtor info after %d attempts", MaxRetries)
}

func fetchValidators(rpcURL string) (config.ValidatorResponse, error) {
	// Make a GET request to the API
	response, err := http.Get(rpcURL)
	if err != nil {
		return config.ValidatorResponse{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return config.ValidatorResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	var data config.ValidatorResponse

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return config.ValidatorResponse{}, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	fmt.Println(data.Pagination.Total)
	return data, nil
}

func findValidatorInfoCustomType(validators []config.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

func fetchDelegationsWithRetry(rpcURL string) (stakingtypes.DelegationResponses, uint64, error) {
	var data stakingtypes.DelegationResponses
	var err error
	var total uint64
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		data, total, err = fetchDelegations(rpcURL)
		if err == nil {
			return data, total, nil
		}
		fmt.Printf("Error fetching delegations info (attempt %d/%d): %v\n", attempt, MaxRetries, err)
		time.Sleep(time.Duration(Backoff.Seconds() * float64(attempt)))
	}
	return stakingtypes.DelegationResponses{}, 0, fmt.Errorf("failed to fetch delegations info after %d attempts", MaxRetries)
}

func fetchDelegations(rpcURL string) (stakingtypes.DelegationResponses, uint64, error) {
	// Make a GET request to the API
	response, err := http.Get(rpcURL)
	if err != nil {
		return nil, 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading response body: %w", err)
	}

	var data config.QueryValidatorDelegationsResponse

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	fmt.Println(data.Pagination.Total)
	total, err := strconv.ParseUint(data.Pagination.Total, 10, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing total from pagination: %w", err)
	}

	return data.DelegationResponses, total, nil
}

func getValidators(stakingClient stakingtypes.QueryClient, blockHeight string) ([]stakingtypes.Validator, error) {
	// Get validator
	ctx := metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight)
	req := &stakingtypes.QueryValidatorsRequest{
		Pagination: &query.PageRequest{
			Limit: LimitPerPage,
		},
	}

	resp, err := stakingClient.Validators(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get validators: %w", err)
	}

	validatorsInfo := resp.Validators
	if validatorsInfo == nil {
		return nil, fmt.Errorf("validators response is nil")
	}

	return validatorsInfo, nil
}
