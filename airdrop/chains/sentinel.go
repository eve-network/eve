package chains

import (
	"fmt"
	"log"
	"strconv"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"
	"github.com/joho/godotenv"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Sentinel() ([]banktypes.Balance, []config.Reward, int, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading Sentinel environment variables: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	blockHeight, err := utils.GetLatestHeight(config.GetSentinelConfig().RPC + "/status")
	if err != nil {
		log.Printf("Failed to get latest height for Sentinel: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Sentinel: %w", err)
	}

	grpcConn, err := utils.SetupGRPCConnection(config.GetSentinelConfig().GRPCAddr)
	if err != nil {
		log.Printf("Failed to connect to gRPC Sentinel: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Sentinel: %w", err)
	}
	defer grpcConn.Close()

	stakingClient := stakingtypes.NewQueryClient(grpcConn)
	validators, err := utils.GetValidators(stakingClient, blockHeight)
	if err != nil {
		log.Printf("Failed to get Sentinel validators: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get Sentinel validators: %w", err)
	}
	log.Println("Validators: ", len(validators))

	delegators := []stakingtypes.DelegationResponse{}
	for validatorIndex, validator := range validators {
		delegationsResponse, err := utils.GetValidatorDelegations(stakingClient, validator.OperatorAddress, blockHeight)
		if err != nil {
			log.Printf("Failed to query delegate info for Sentinel validator: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to query delegate info for Sentinel validator: %w", err)
		}
		log.Println("Response ", len(delegationsResponse.DelegationResponses))
		log.Println("Sentinel validator "+strconv.Itoa(validatorIndex)+" ", delegationsResponse.Pagination.Total)
		delegators = append(delegators, delegationsResponse.DelegationResponses...)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")
	tokenInUsd, err := utils.FetchTokenPrice(config.GetSentinelConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Sentinel token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Sentinel token price: %w", err)
	}
	tokenIn20Usd := usd.Quo(tokenInUsd)

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}
	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")

	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}

	eveAirdrop, err := sdkmath.LegacyNewDecFromStr(config.EveAirdrop)
	if err != nil {
		log.Printf("Failed to convert EveAirdrop string to dec: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to convert EveAirdrop string to dec: %w", err)
	}

	testAmount := sdkmath.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}

		eveAirdropAmount := (eveAirdrop.MulInt64(int64(config.GetSentinelConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Sentinel bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		reward := config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdropAmount,
			ChainID:         config.GetSentinelConfig().ChainID,
		}
		rewardInfo = append(rewardInfo, reward)

		balance := banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdropAmount.TruncateInt())),
		}
		balanceInfo = append(balanceInfo, balance)
		testAmount = eveAirdropAmount.Add(testAmount)
	}

	log.Println("Sentinel balance: ", testAmount)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
