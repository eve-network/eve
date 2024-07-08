package chains

import (
	"fmt"
	"log"
	"strconv"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Terra() ([]banktypes.Balance, []config.Reward, int, error) {
	var delegators []stakingtypes.DelegationResponse

	validatorsResponse, err := utils.FetchValidators(config.GetTerraConfig().API, config.LimitPerPage)
	if err != nil {
		log.Printf("Failed to fetch validators for Terra: %v\n", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch validators for Terra: %w", err)
	}

	validators := validatorsResponse.Validators
	log.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		url := fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations?pagination.limit=%d&pagination.count_total=true",
			config.GetTerraConfig().API, validator.OperatorAddress, config.LimitPerPage)
		delegations, total, err := utils.FetchDelegations(url)
		if err != nil {
			log.Printf("Failed to fetch delegations for validator %s at Terra using URL %s: %v\n", validator.OperatorAddress, url, err)
			return nil, nil, 0, fmt.Errorf("failed to fetch delegations for Terra: %w", err)
		}
		log.Println(validator.OperatorAddress)
		log.Println("Response ", len(delegations))
		log.Println("Terra validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegations...)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")

	tokenInUsd, err := utils.FetchTokenPrice(config.GetTerraConfig().CoinID)
	if err != nil {
		log.Printf("Failed to fetch Terra token price: %v\n", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Terra token price: %w", err)
	}
	tokenIn20Usd := usd.Quo(tokenInUsd)

	var rewardInfo []config.Reward
	var balanceInfo []banktypes.Balance

	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}

	eveAirdrop, err := sdkmath.LegacyNewDecFromStr(config.EveAirdrop)
	if err != nil {
		log.Printf("Failed to convert EveAirdrop string to dec: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to convert string to dec: %w", err)
	}

	testAmount := sdkmath.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}

		if totalTokenDelegate.IsZero() {
			return nil, nil, 0, fmt.Errorf("total token delegate is zero, cannot proceed with airdrop calculation")
		}

		eveReward := eveAirdrop.MulInt64(int64(config.GetTerraConfig().Percent)).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			log.Printf("Failed to convert Terra bech32 address for delegator %s: %v\n", delegator.Delegation.DelegatorAddress, err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveReward,
			ChainID:         config.GetTerraConfig().ChainID,
		})

		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveReward.TruncateInt())),
		})

		testAmount = eveReward.Add(testAmount)
	}

	log.Println("Terra balance: ", testAmount)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
