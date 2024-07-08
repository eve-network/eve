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
	"github.com/ethereum/go-ethereum/common"
	airdropBackoff "github.com/eve-network/eve/airdrop/backoff"
	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"cosmossdk.io/core/address"
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// handleHTTPError handles HTTP-related errors by wrapping them with a message.
func handleHTTPError(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}
	return nil
}

// MakeGetRequest creates and sends a GET request to the specified URI.
func MakeGetRequest(uri string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, handleHTTPError(err, "failed to create HTTP request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, handleHTTPError(err, "failed to send HTTP request")
	}

	return res, nil
}

// GetFunctionName returns the name of the given function.
func GetFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// FindValidatorInfo finds the index of a validator by its address.
func FindValidatorInfo(validators []stakingtypes.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

// GetLatestHeight fetches the latest block height from the given API URL.
func GetLatestHeight(apiURL string) (string, error) {
	ctx := context.Background()
	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	var response *http.Response
	var err error

	retryableRequest := func() error {
		response, err = MakeGetRequest(apiURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return "", handleHTTPError(err, "error making GET request to get latest height")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", handleHTTPError(err, "error reading response body")
	}

	var data config.NodeResponse
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", handleHTTPError(err, "error unmarshalling JSON")
	}

	latestBlockHeight := data.Result.SyncInfo.LatestBlockHeight
	log.Println("Block height:", latestBlockHeight)

	return latestBlockHeight, nil
}

// GetValidators retrieves the list of validators at a specific block height.
func GetValidators(stakingClient stakingtypes.QueryClient, blockHeight string) ([]stakingtypes.Validator, error) {
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

// GetValidatorDelegations retrieves the delegations for a specific validator at a specific block height.
func GetValidatorDelegations(stakingClient stakingtypes.QueryClient, validatorAddr, blockHeight string) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
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

// ConvertBech32Address converts an address from another chain to a Bech32 address.
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

// FindValidatorInfoCustomType finds the index of a validator by its address in a custom validator type.
func FindValidatorInfoCustomType(validators []config.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

// FetchValidators fetches the list of validators from the given RPC URL.
func FetchValidators(chainAPI string, limitPerPage int) (config.ValidatorResponse, error) {
	ctx := context.Background()
	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	var response *http.Response
	var err error

	rpcURL := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators?pagination.limit=%d&pagination.count_total=true",
		chainAPI, limitPerPage)

	retryableRequest := func() error {
		response, err = MakeGetRequest(rpcURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return config.ValidatorResponse{}, fmt.Errorf("error making GET request to fetch validators: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return config.ValidatorResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	var data config.ValidatorResponse
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return config.ValidatorResponse{}, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	log.Println(data.Pagination.Total)
	return data, nil
}

// FetchDelegations fetches delegations from the given RPC URL.
func FetchDelegations(rpcURL string) (stakingtypes.DelegationResponses, uint64, error) {
	ctx := context.Background()
	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	var response *http.Response
	var err error

	retryableRequest := func() error {
		response, err = MakeGetRequest(rpcURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return nil, 0, fmt.Errorf("error making GET request to fetch delegations: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading response body: %w", err)
	}

	var data config.QueryValidatorDelegationsResponse
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return nil, 0, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	log.Println(data.Pagination.Total)
	total, err := strconv.ParseUint(data.Pagination.Total, 10, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing total from pagination: %w", err)
	}

	return data.DelegationResponses, total, nil
}

// BalanceFunction defines a function type that returns balance info, reward info, and length.
type BalanceFunction func() ([]banktypes.Balance, []config.Reward, int, error)

// RetryableBalanceFunc wraps a BalanceFunction with retry logic.
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

// FetchTokenPrice fetches the token price from the given API URL.
func FetchTokenPrice(coinID string) (sdkmath.LegacyDec, error) {
	ctx := context.Background()
	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	var response *http.Response
	var err error

	apiURL := fmt.Sprintf("%s%s&vs_currencies=usd", config.APICoingecko, coinID)

	retryableRequest := func() error {
		response, err = MakeGetRequest(apiURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error making GET request to fetch token price: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error reading response body for token price: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
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

// SetupGRPCConnection sets up a gRPC connection to the specified address.
func SetupGRPCConnection(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// FetchTokenInfo fetches token information from the specified contract address.
func FetchTokenInfo(token, contractAddress, apiFromConfig string) (config.NftHolder, error) {
	queryString := fmt.Sprintf(`{"all_nft_info":{"token_id":%s}}`, token)
	encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
	apiURL := fmt.Sprintf("%s/cosmwasm/wasm/v1/contract/%s/smart/%s", apiFromConfig, contractAddress, encodedQuery)

	ctx := context.Background()
	exponentialBackoff := airdropBackoff.NewBackoff(ctx)

	var response *http.Response
	var err error

	retryableRequest := func() error {
		response, err = MakeGetRequest(apiURL)
		return err
	}

	if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
		return config.NftHolder{}, fmt.Errorf("error making GET request to fetch token info: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return config.NftHolder{}, fmt.Errorf("error reading response body when fetching token info: %w", err)
	}

	var data config.TokenInfoResponse
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return config.NftHolder{}, fmt.Errorf("error unmarshalling JSON when fetching token info: %w", err)
	}

	return config.NftHolder{
		Address: data.Data.Access.Owner,
		TokenID: token,
	}, nil
}

// FetchTokenIds fetches all token IDs from the specified contract address.
func FetchTokenIds(contractAddress, apiFromConfig string) ([]string, error) {
	paginationKey := "0"
	tokenIds := []string{}

	for {
		queryString := fmt.Sprintf(`{"all_tokens":{"limit":1000,"start_after":"%s"}}`, paginationKey)
		encodedQuery := base64.StdEncoding.EncodeToString([]byte(queryString))
		apiURL := fmt.Sprintf("%s/cosmwasm/wasm/v1/contract/%s/smart/%s", apiFromConfig, contractAddress, encodedQuery)

		ctx := context.Background()
		exponentialBackoff := airdropBackoff.NewBackoff(ctx)

		var response *http.Response
		var err error

		retryableRequest := func() error {
			response, err = MakeGetRequest(apiURL)
			return err
		}

		if err := backoff.Retry(retryableRequest, exponentialBackoff); err != nil {
			return nil, fmt.Errorf("error making GET request to fetch token ids: %w", err)
		}
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body when fetching token ids: %w", err)
		}

		var data config.TokenIdsResponse
		if err := json.Unmarshal(responseBody, &data); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON when fetching token ids: %w", err)
		}

		tokenIds = append(tokenIds, data.Data.Token...)
		if len(data.Data.Token) == 0 {
			break
		}

		paginationKey = data.Data.Token[len(data.Data.Token)-1]
		log.Println("pagination key:", paginationKey)
		if len(paginationKey) == 0 {
			break
		}
	}

	log.Println(len(tokenIds))
	return tokenIds, nil
}

func StringFromEthAddress(codec address.Codec, ethAddress common.Address) (string, error) {
	return codec.BytesToString(ethAddress.Bytes())
}
