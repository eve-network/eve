package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strconv"

	"github.com/cenkalti/backoff/v4"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	airdropBackoff "github.com/eve-network/eve/airdrop/backoff"
	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc/metadata"
)

const (
	LimitPerPage = 100000000
)

func MakeGetRequest(uri string) (*http.Response, error) {
	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Send the HTTP request and get the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
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
		return "", fmt.Errorf("error making GET request to get latest height: %w", err)
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

func GetValidators(stakingClient stakingtypes.QueryClient, blockHeight string) ([]stakingtypes.Validator, error) {
	// Get validator
	ctx := metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight)
	req := &stakingtypes.QueryValidatorsRequest{
		Pagination: &query.PageRequest{
			Limit: LimitPerPage,
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

func GetValidatorDelegations(stakingClient stakingtypes.QueryClient, validatorAddr string, blockHeight string) (
	*stakingtypes.QueryValidatorDelegationsResponse, error,
) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight)
	req := &stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: validatorAddr,
		Pagination: &query.PageRequest{
			CountTotal: true,
			Limit:      LimitPerPage,
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

	fmt.Println(data.Pagination.Total)
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
		return nil, 0, fmt.Errorf("error making GET request to get fetch delegations: %w", err)
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