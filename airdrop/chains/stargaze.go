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

func Stargaze() ([]banktypes.Balance, []config.Reward, int, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading Stargaze environment variables: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	blockHeight, err := utils.GetLatestHeight(config.GetStargazeConfig().RPC + "/status")
	if err != nil {
		log.Printf("Failed to get latest height for Stargaze: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Stargaze: %w", err)
	}

	grpcConn, err := utils.SetupGRPCConnection(config.GetStargazeConfig().GRPCAddr)
	if err != nil {
		log.Printf("Failed to connect to gRPC Stargaze: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Stargaze: %w", err)
	}
	defer grpcConn.Close()

	stakingClient := stakingtypes.NewQueryClient(grpcConn)
	validators, err := utils.GetValidators(stakingClient, blockHeight)
	if err != nil {
		log.Printf("Failed to get Stargaze validators: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get Stargaze validators: %w", err)
	}

	var delegators []stakingtypes.DelegationResponse
	for validatorIndex, validator := range validators {
		delegationsResponse, err := utils.GetValidatorDelegations(stakingClient, validator.OperatorAddress, blockHeight)
		if err != nil {
			log.Printf("Failed to query delegate info for Stagaze validator: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to query delegate info for Stargaze validator: %w", err)
		}
		log.Printf("Validator %d: %d delegations", validatorIndex, len(delegationsResponse.DelegationResponses))
		delegators = append(delegators, delegationsResponse.DelegationResponses...)
	}

	// Calculate token price and threshold
	minimumTokensThreshold, err := utils.GetMinimumTokensThreshold(config.GetStargazeConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Stargaze token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Stargaze token price: %w", err)
	}

	var rewardInfo []config.Reward
	var balanceInfo []banktypes.Balance
	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		token := (delegator.Delegation.Shares.MulInt(validators[validatorIndex].Tokens)).QuoTruncate(validators[validatorIndex].DelegatorShares)
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
		token := (delegator.Delegation.Shares.MulInt(validators[validatorIndex].Tokens)).QuoTruncate(validators[validatorIndex].DelegatorShares)
		if token.LT(minimumTokensThreshold) {
			continue
		}

		eveReward := eveAirdrop.MulInt64(int64(config.GetStargazeConfig().Percent)).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Stargaze bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveReward,
			ChainID:         config.GetStargazeConfig().ChainID,
		})

		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveReward.TruncateInt())),
		})

		testAmount = eveReward.Add(testAmount)
	}

	log.Printf("Stargaze balance: %s", testAmount)

	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
