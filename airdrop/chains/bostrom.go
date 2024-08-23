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

func Bostrom() ([]banktypes.Balance, []config.Reward, int, error) {
	delegators := []stakingtypes.DelegationResponse{}

	// Fetch validators
	validatorsResponse, err := utils.FetchValidators(config.GetBostromConfig().API, config.LimitPerPage)
	if err != nil {
		log.Printf("Failed to fetch Bostrom validators: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch validators for Bostrom: %w", err)
	}
	validators := validatorsResponse.Validators
	log.Println("Validators: ", len(validators))

	// Fetch delegations for each validator
	for validatorIndex, validator := range validators {
		delegations, total, err := utils.FetchDelegations(
			config.GetBostromConfig().API,
			validator.OperatorAddress,
			config.LimitPerPage,
		)
		if err != nil {
			log.Printf("Failed to fetch delegations for Bostrom: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to fetch delegations for Bostrom: %w", err)
		}
		log.Printf("Bostrom validator %d has %d delegations", validatorIndex, total)
		delegators = append(delegators, delegations...)
	}

	// Calculate token price and threshold
	minimumTokensThreshold, err := utils.GetMinimumTokensThreshold(config.GetBostromConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Bostrom token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Bostrom token price: %w", err)
	}

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}
	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")

	// Calculate total tokens delegated
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
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

	// Calculate rewards and balances for delegators
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(minimumTokensThreshold) {
			continue
		}

		eveAirdropTokens := (eveAirdrop.MulInt64(int64(config.GetBostromConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Bostrom bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdropTokens,
			ChainID:         config.GetBostromConfig().ChainID,
		})
		testAmount = eveAirdropTokens.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdropTokens.TruncateInt())),
		})
	}

	log.Println("Bostrom balance: ", testAmount)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
