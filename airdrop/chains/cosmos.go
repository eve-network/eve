package chains

import (
	"fmt"
	"log"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Cosmos() ([]banktypes.Balance, []config.Reward, int, error) {
	var delegators []stakingtypes.DelegationResponse

	// Fetch validators
	validatorsResponse, err := utils.FetchValidators(config.GetCosmosHubConfig().API, config.LimitPerPage)
	if err != nil {
		log.Printf("Failed to fetch validator for Cosmos: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch validator for Cosmos: %w", err)
	}
	validators := validatorsResponse.Validators
	log.Printf("Validators: %d", len(validators))

	// Fetch delegations for each validator
	for validatorIndex, validator := range validators {
		delegations, total, err := utils.FetchDelegations(
			config.GetCosmosHubConfig().API,
			validator.OperatorAddress,
			config.LimitPerPage,
		)
		if err != nil {
			log.Printf("Failed to fetch Delegations for Cosmos: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to fetch Delegations for Cosmos: %w", err)
		}
		log.Printf("Cosmos validator %d has %d delegators", validatorIndex, total)
		log.Printf("Validator %s: %d delegations (total: %d)", validator.OperatorAddress, len(delegations), total)
		delegators = append(delegators, delegations...)
	}

	// Calculate token price and threshold
	minimumTokensThreshold, err := utils.GetMinimumTokensThreshold(config.GetCosmosHubConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Cosmos token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Cosmos token price: %w", err)
	}

	var (
		rewardInfo         []config.Reward
		balanceInfo        []banktypes.Balance
		totalTokenDelegate = sdkmath.LegacyMustNewDecFromStr("0")
	)

	// Calculate total delegated tokens
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := delegator.Delegation.Shares.MulInt(validatorInfo.Tokens).QuoTruncate(validatorInfo.DelegatorShares)
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}

	eveAirdrop, err := sdkmath.LegacyNewDecFromStr(config.EveAirdrop)
	if err != nil {
		log.Printf("Failed to convert EveAirdrop string to dec: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to convert EveAirdrop string to dec: %w", err)
	}

	totalEveAirdrop := sdkmath.LegacyMustNewDecFromStr("0")

	// Process each delegator
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := delegator.Delegation.Shares.MulInt(validatorInfo.Tokens).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(minimumTokensThreshold) {
			continue
		}

		eveAirdropToken := eveAirdrop.MulInt64(int64(config.GetCosmosHubConfig().Percent)).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Cosmos bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdropToken,
			ChainID:         config.GetCosmosHubConfig().ChainID,
		})

		totalEveAirdrop = totalEveAirdrop.Add(eveAirdropToken)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdropToken.TruncateInt())),
		})
	}

	log.Printf("Cosmos balance: %s", totalEveAirdrop)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
