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

func Composable() ([]banktypes.Balance, []config.Reward, int, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading Composable environment variables: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	// Get the latest block height
	blockHeight, err := utils.GetLatestHeight(config.GetComposableConfig().RPC + "/status")
	if err != nil {
		log.Printf("Failed to get latest height for Composable: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Composable: %w", err)
	}

	// Setup gRPC connection
	grpcConn, err := utils.SetupGRPCConnection(config.GetComposableConfig().GRPCAddr)
	if err != nil {
		log.Printf("Failed to connect to gRPC Composable: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Composable: %w", err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	// Fetch validators
	validators, err := utils.GetValidators(stakingClient, blockHeight)
	if err != nil {
		log.Printf("Failed to get Composable validators: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get Composable validators: %w", err)
	}
	log.Println("Validators: ", len(validators))

	// Fetch delegations for each validator
	var delegators []stakingtypes.DelegationResponse
	for validatorIndex, validator := range validators {
		delegationsResponse, err := utils.GetValidatorDelegations(stakingClient, validator.OperatorAddress, blockHeight)
		if err != nil {
			log.Printf("Failed to query delegate info for Composable validator: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to query delegate info for Composable validator: %w", err)
		}
		total := delegationsResponse.Pagination.Total
		log.Printf("Composable validator %d has %d delegators", validatorIndex, total)
		log.Printf("Validator %s: %d delegations (total: %d)", validator.OperatorAddress, len(delegationsResponse.DelegationResponses), total)
		delegators = append(delegators, delegationsResponse.DelegationResponses...)
	}

	// Calculate token price and threshold
	minimumTokensThreshold, err := utils.GetMinimumTokensThreshold(config.GetComposableConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Composable token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Composable token price: %w", err)
	}

	// Process delegations and calculate rewards
	var (
		rewardInfo         []config.Reward
		balanceInfo        []banktypes.Balance
		totalTokenDelegate = sdkmath.LegacyMustNewDecFromStr("0")
	)

	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := delegator.Delegation.Shares.MulInt(validatorInfo.Tokens).QuoTruncate(validatorInfo.DelegatorShares)
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
		token := delegator.Delegation.Shares.MulInt(validatorInfo.Tokens).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(minimumTokensThreshold) {
			continue
		}
		eveAirdropToken := eveAirdrop.MulInt64(int64(config.GetComposableConfig().Percent)).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Composable bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdropToken,
			ChainID:         config.GetComposableConfig().ChainID,
		})

		testAmount = testAmount.Add(eveAirdropToken)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdropToken.TruncateInt())),
		})
	}

	log.Println("Composable balance: ", testAmount)

	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
