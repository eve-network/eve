package chains

import (
	"fmt"
	"log"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"
	"github.com/joho/godotenv"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Akash() ([]banktypes.Balance, []config.Reward, int, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading Akash environment variables: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	// Get latest block height
	blockHeight, err := utils.GetLatestHeight(config.GetAkashConfig().RPC + "/status")
	if err != nil {
		log.Printf("Failed to get latest height for Akash: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Akash: %w", err)
	}

	// Setup gRPC connection
	grpcAddr := config.GetAkashConfig().GRPCAddr
	grpcConn, err := utils.SetupGRPCConnection(grpcAddr)
	if err != nil {
		log.Printf("Failed to connect to gRPC Akash: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Akash: %w", err)
	}
	defer grpcConn.Close()

	stakingClient := stakingtypes.NewQueryClient(grpcConn)
	delegators := []stakingtypes.DelegationResponse{}

	// Get validators
	validators, err := utils.GetValidators(stakingClient, blockHeight)
	if err != nil {
		log.Printf("Failed to get Akash validators: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get Akash validators: %w", err)
	}
	log.Println("Validators: ", len(validators))

	// Get delegations for each validator
	for validatorIndex, validator := range validators {
		delegationsResponse, err := utils.GetValidatorDelegations(stakingClient, validator.OperatorAddress, blockHeight)
		if err != nil {
			log.Printf("Failed to query delegate info for Akash validator: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to query delegate info for Akash validator: %w", err)
		}
		total := delegationsResponse.Pagination.Total
		log.Printf("Akash validator %d has %d delegators", validatorIndex, total)
		delegators = append(delegators, delegationsResponse.DelegationResponses...)
	}

	// Calculate token price and threshold
	minimumTokensThreshold, err := utils.GetMinimumTokensThreshold(config.GetAkashConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Akash token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Akash token price: %w", err)
	}

	// Prepare for airdrop calculation
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

		if token.LT(minimumTokensThreshold) {
			continue
		}

		eveAirdropTokens := (eveAirdrop.MulInt64(int64(config.GetAkashConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Akash bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdropTokens,
			ChainID:         config.GetAkashConfig().ChainID,
		})
		testAmount = eveAirdropTokens.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdropTokens.TruncateInt())),
		})
	}

	log.Println("Akash balance: ", testAmount)

	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
