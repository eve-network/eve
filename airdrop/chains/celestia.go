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

func Celestia() ([]banktypes.Balance, []config.Reward, int, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading Celestial environment variables: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to load env: %w", err)
	}

	// Get the latest block height
	blockHeight, err := utils.GetLatestHeight(config.GetCelestiaConfig().RPC + "/status")
	if err != nil {
		log.Printf("Failed to get latest height for Celestial: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Celestia: %w", err)
	}

	// Setup gRPC connection
	grpcConn, err := utils.SetupGRPCConnection(config.GetCelestiaConfig().GRPCAddr)
	if err != nil {
		log.Printf("Failed to connect to gRPC Celestial: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to connect to gRPC Celestia: %w", err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	// Fetch validators
	validators, err := utils.GetValidators(stakingClient, blockHeight)
	if err != nil {
		log.Printf("Failed to get Celestial validators: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get Celestia validators: %w", err)
	}
	log.Println("Validators: ", len(validators))

	// Fetch delegations for each validator
	var delegators []stakingtypes.DelegationResponse
	for validatorIndex, validator := range validators {
		url := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations?pagination.limit=%d&pagination.count_total=true",
			config.GetCelestiaConfig().API, validator.OperatorAddress, config.LimitPerPage)
		delegations, total, err := utils.FetchDelegations(url)
		if err != nil {
			log.Printf("Failed to query delegate info for Celestial validator: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to fetch delegations for Celestia: %w", err)
		}
		log.Println("Validator:", validator.OperatorAddress, "Index:", validatorIndex, "Total:", total)
		delegators = append(delegators, delegations...)
	}

	// Fetch token price in USD
	tokenInUsd, err := utils.FetchTokenPrice(config.GetCelestiaConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Celestial token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Celestia token price: %w", err)
	}
	tokenIn20Usd := sdkmath.LegacyMustNewDecFromStr("20").Quo(tokenInUsd)

	// Process delegations and calculate rewards
	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")
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

	var rewardInfo []config.Reward
	var balanceInfo []banktypes.Balance
	testAmount := sdkmath.LegacyMustNewDecFromStr("0")

	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := delegator.Delegation.Shares.MulInt(validatorInfo.Tokens).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdropToken := eveAirdrop.MulInt64(int64(config.GetCelestiaConfig().Percent)).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Celestial bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdropToken,
			ChainID:         config.GetCelestiaConfig().ChainID,
		})

		testAmount = testAmount.Add(eveAirdropToken)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdropToken.TruncateInt())),
		})
	}

	log.Println("Celestia balance: ", testAmount)

	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
