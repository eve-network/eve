package main

// error max size response
import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func celestia() ([]banktypes.Balance, []config.Reward, int, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	blockHeight, err := getLatestHeightWithRetry(config.GetCelestiaConfig().RPC + "/status")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Celestia: %w", err)
	}

	grpcAddr := config.GetCelestiaConfig().GRPCAddr
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Celestia: %w", err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	delegators := []stakingtypes.DelegationResponse{}

	validators, err := getValidators(stakingClient, blockHeight)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get Celestia validators: %w", err)
	}

	fmt.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		url := config.GetCelestiaConfig().API + "/cosmos/staking/v1beta1/validators/" + validator.OperatorAddress + "/delegations?pagination.limit=" + strconv.Itoa(LimitPerPage) + "&pagination.count_total=true"
		delegations, total, err := fetchDelegationsWithRetry(url)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to fetch delegations for Celestia: %w", err)
		}
		fmt.Println(validator.OperatorAddress)
		fmt.Println("Response ", len(delegations))
		fmt.Println("Celestia validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegations...)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")

	apiURL := APICoingecko + config.GetCelestiaConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd, err := fetchCelestiaTokenPriceWithRetry(apiURL)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch Celestia token price: %w", err)
	}
	tokenIn20Usd := usd.QuoTruncate(tokenInUsd)

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}
	eveAirdrop := sdkmath.LegacyMustNewDecFromStr(EveAirdrop)
	testAmount, _ := sdkmath.LegacyNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetCelestiaConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := convertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetCelestiaConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println("Celestia ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}

func fetchCelestiaTokenPriceWithRetry(apiURL string) (sdkmath.LegacyDec, error) {
	var data sdkmath.LegacyDec
	var err error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		data, err = fetchCelestiaTokenPrice(apiURL)
		if err == nil {
			return data, nil
		}

		fmt.Printf("error fetching Celestia token price (attempt %d/%d): %v\n", attempt, MaxRetries, err)

		if attempt < MaxRetries {
			// Calculate backoff duration using exponential backoff strategy
			backoffDuration := time.Duration(Backoff.Seconds() * math.Pow(2, float64(attempt)))
			fmt.Printf("retrying after %s...\n", backoffDuration)
			time.Sleep(backoffDuration)
		}
	}

	return sdkmath.LegacyDec{}, fmt.Errorf("failed to fetch Celestia token price after %d attempts: %v", MaxRetries, err)
}

func fetchCelestiaTokenPrice(apiURL string) (sdkmath.LegacyDec, error) {
	// Make a GET request to the API
	response, err := makeGetRequest(apiURL)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error making GET request to fetch Celestia token price: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error reading response body for Celestia token price: %w", err)
	}

	var data config.CelestiaPrice

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error unmarshalling JSON for Celestia token price: %w", err)
	}

	tokenInUsd := sdkmath.LegacyMustNewDecFromStr(data.Token.USD.String())
	return tokenInUsd, nil
}
