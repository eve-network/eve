package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func akash() ([]banktypes.Balance, []config.Reward, int, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	blockHeight, err := getLatestHeightWithRetry(config.GetAkashConfig().RPC + "/status")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Akash: %w", err)
	}

	grpcAddr := config.GetAkashConfig().GRPCAddr
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Akash: %w", err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	delegators := []stakingtypes.DelegationResponse{}

	validators, err := getValidators(stakingClient, blockHeight)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get Akash validators: %w", err)
	}
	fmt.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		var header metadata.MD
		delegationsResponse, _ := stakingClient.ValidatorDelegations(
			metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight), // Add metadata to request
			&stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: validator.OperatorAddress,
				Pagination: &query.PageRequest{
					CountTotal: true,
					Limit:      LimitPerPage,
				},
			},
			grpc.Header(&header), // Retrieve header from response
		)
		total := delegationsResponse.Pagination.Total
		fmt.Println("Response ", len(delegationsResponse.DelegationResponses))
		fmt.Println("Akash validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegationsResponse.DelegationResponses...)
	}

	usd := math.LegacyMustNewDecFromStr("20")

	apiURL := APICoingecko + config.GetAkashConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd, err := fetchAkashTokenPriceWithRetry(apiURL)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch Akash token price: %w", err)
	}
	tokenIn20Usd := usd.QuoTruncate(tokenInUsd)

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenDelegate := math.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}
	eveAirdrop := math.LegacyMustNewDecFromStr(EveAirdrop)
	testAmount, err := math.LegacyNewDecFromStr("0")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to convert string to dec: %w", err)
	}
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetAkashConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address := convertBech32Address(delegator.Delegation.DelegatorAddress)
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetAkashConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println("Akash ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}

func fetchAkashTokenPriceWithRetry(apiURL string) (math.LegacyDec, error) {
	var data math.LegacyDec
	var err error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		data, err = fetchAkashTokenPrice(apiURL)
		if err == nil {
			return data, nil
		}
		fmt.Printf("error fetching Akash token price (attempt %d/%d): %v\n", attempt, MaxRetries, err)
		time.Sleep(time.Duration(time.Duration(attempt * Backoff).Milliseconds()))
	}
	return math.LegacyDec{}, fmt.Errorf("failed to fetch Akash token price after %d attempts", MaxRetries)
}

func fetchAkashTokenPrice(apiURL string) (math.LegacyDec, error) {
	// Make a GET request to the API
	response, err := http.Get(apiURL)
	if err != nil {
		return math.LegacyDec{}, fmt.Errorf("error making GET request to fetch Akash token price: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return math.LegacyDec{}, fmt.Errorf("error reading response body for Akash token price: %w", err)
	}

	var data config.AkashPrice

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return math.LegacyDec{}, fmt.Errorf("error unmarshalling JSON for Akash token price: %w", err)
	}

	tokenInUsd := math.LegacyMustNewDecFromStr(data.Token.USD.String())
	return tokenInUsd, nil
}
