package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/cenkalti/backoff/v4"
	airdropBackoff "github.com/eve-network/eve/airdrop/backoff"
	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func handleHTTPError(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}
	return nil
}

func MakeGetRequest(uri string) (*http.Response, error) {
	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, handleHTTPError(err, "failed to create HTTP request")
	}

	// Send the HTTP request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, handleHTTPError(err, "failed to send HTTP request")
	}

	return res, nil
}

func GetFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func FindValidatorInfo(validators []stakingtypes.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

func GetLatestHeight(apiURL string) (string, error) {
	ctx := context.Background()

	var response *http.Response
	var err error

	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	retryableRequest := func() error {
		// Make a GET request to the API
		response, err = MakeGetRequest(apiURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return "", handleHTTPError(err, "error making GET request to get latest height")
	}

	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", handleHTTPError(err, "error reading response body")
	}

	// Parse the response body into a NodeResponse struct
	var data config.NodeResponse
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", handleHTTPError(err, "error unmarshalling JSON")
	}

	// Extract the latest block height from the response
	latestBlockHeight := data.Result.SyncInfo.LatestBlockHeight
	log.Println("Block height:", latestBlockHeight)

	return latestBlockHeight, nil
}

func GetValidators(stakingClient stakingtypes.QueryClient, blockHeight string) ([]stakingtypes.Validator, error) {
	// Get validator
	ctx := metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight)
	req := &stakingtypes.QueryValidatorsRequest{
		Pagination: &query.PageRequest{
			Limit: config.LimitPerPage,
		},
	}

	var resp *stakingtypes.QueryValidatorsResponse
	var err error

	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	retryableRequest := func() error {
		resp, err = stakingClient.Validators(ctx, req)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return nil, fmt.Errorf("failed to get validators: %w", err)
	}

	if resp == nil || resp.Validators == nil {
		return nil, fmt.Errorf("validators response is nil")
	}

	return resp.Validators, nil
}

func GetValidatorDelegations(stakingClient stakingtypes.QueryClient, validatorAddr, blockHeight string) (
	*stakingtypes.QueryValidatorDelegationsResponse, error,
) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight)
	req := &stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: validatorAddr,
		Pagination: &query.PageRequest{
			CountTotal: true,
			Limit:      config.LimitPerPage,
		},
	}

	var resp *stakingtypes.QueryValidatorDelegationsResponse
	var err error

	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	retryableRequest := func() error {
		resp, err = stakingClient.ValidatorDelegations(ctx, req)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return nil, fmt.Errorf("failed to get validator delegations: %w", err)
	}

	return resp, nil
}

func ConvertBech32Address(otherChainAddress string) (string, error) {
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

func FindValidatorInfoCustomType(validators []config.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

func FetchValidators(rpcURL string) (config.ValidatorResponse, error) {
	ctx := context.Background()

	var response *http.Response
	var err error

	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	retryableRequest := func() error {
		// Make a GET request to the API
		response, err = MakeGetRequest(rpcURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return config.ValidatorResponse{}, fmt.Errorf("error making GET request to get fetch validators: %w", err)
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

	log.Println(data.Pagination.Total)
	return data, nil
}

func FetchDelegations(rpcURL string) (stakingtypes.DelegationResponses, uint64, error) {
	ctx := context.Background()

	var response *http.Response
	var err error

	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	retryableRequest := func() error {
		// Make a GET request to the API
		response, err = MakeGetRequest(rpcURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return nil, 0, fmt.Errorf("error making GET request to fetch delegations: %w", err)
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

	log.Println(data.Pagination.Total)
	total, err := strconv.ParseUint(data.Pagination.Total, 10, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing total from pagination: %w", err)
	}

	return data.DelegationResponses, total, nil
}

// Define a function type that returns balance info, reward info and length
type BalanceFunction func() ([]banktypes.Balance, []config.Reward, int, error)

// Retryable function to wrap balanceFunction with retry logic
func RetryableBalanceFunc(fn BalanceFunction) BalanceFunction {
	return func() ([]banktypes.Balance, []config.Reward, int, error) {
		for attempt := 1; attempt <= config.MaxRetries; attempt++ {
			balances, rewards, length, err := fn()
			if err == nil {
				return balances, rewards, length, nil
			}
			fmt.Printf("Failed attempt %d for function %s: %v\n", attempt, GetFunctionName(fn), err)
		}
		return nil, nil, 0, fmt.Errorf("maximum retries reached for function %s", GetFunctionName(fn))
	}
}

func FetchTokenPrice(apiURL, coinID string) (sdkmath.LegacyDec, error) {
	ctx := context.Background()

	var response *http.Response
	var err error

	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	retryableRequest := func() error {
		// Make a GET request to the API
		response, err = MakeGetRequest(apiURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error making GET request to fetch token price: %w", err)
	}

	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error reading response body for token price: %w", err)
	}

	var data interface{}
	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error unmarshalling JSON for token price: %w", err)
	}
	tokenPrice := data.(map[string]interface{})
	priceInUsd := fmt.Sprintf("%v", tokenPrice[coinID].(map[string]interface{})["usd"])
	var tokenPriceInUsd sdkmath.LegacyDec

	if strings.Contains(priceInUsd, "e-") {
		rawPrice := strings.Split(priceInUsd, "e-")
		base := rawPrice[0]
		power := rawPrice[1]
		powerInt, _ := strconv.ParseUint(power, 10, 64)
		baseDec, _ := sdkmath.LegacyNewDecFromStr(base)
		tenDec, _ := sdkmath.LegacyNewDecFromStr("10")
		tokenPriceInUsd = baseDec.Quo(tenDec.Power(powerInt))
	} else {
		tokenPriceInUsd = sdkmath.LegacyMustNewDecFromStr(priceInUsd)
	}
	return tokenPriceInUsd, nil
}

func SetupGRPCConnection(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func FetchTokenInfo(token, contractAddress, apiFromConfig string) (config.NftHolder, error) {
	queryString := fmt.Sprintf(`{"all_nft_info":{"token_id":%s}}`, token)
	encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
	apiURL := apiFromConfig + "/cosmwasm/wasm/v1/contract/" + contractAddress + "/smart/" + encodedQuery

	ctx := context.Background()

	var response *http.Response
	var err error

	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	retryableRequest := func() error {
		// Make a GET request to the API
		response, err = MakeGetRequest(apiURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
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

func FetchTokenIds(contractAddress, apiFromConfig string) ([]string, error) {
	// Make a GET request to the API
	paginationKey := "0"
	index := 0
	tokenIds := []string{}
	for {
		index++
		queryString := fmt.Sprintf(`{"all_tokens":{"limit":1000,"start_after":"%s"}}`, paginationKey)
		encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
		apiURL := apiFromConfig + "/cosmwasm/wasm/v1/contract/" + contractAddress + "/smart/" + encodedQuery
		ctx := context.Background()

		var response *http.Response
		var err error

		exponentialBackoff := airdropBackoff.NewBackoff(ctx)

		retryableRequest := func() error {
			// Make a GET request to the API
			response, err = MakeGetRequest(apiURL)
			return err
		}

		if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
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
			log.Println("pagination key:", paginationKey)
			if len(paginationKey) == 0 {
				break
			}
		}
	}

	log.Println(len(tokenIds))
	return tokenIds, nil
}
