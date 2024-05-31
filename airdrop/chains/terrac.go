package chains

import (
	"fmt"
	"strconv"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Terrac() ([]banktypes.Balance, []config.Reward, int, error) {
	delegators := []stakingtypes.DelegationResponse{}

	rpc := config.GetTerracConfig().API + "/cosmos/staking/v1beta1/validators?pagination.limit=" + strconv.Itoa(config.LimitPerPage) + "&pagination.count_total=true"
	validatorsResponse, err := utils.FetchValidators(rpc)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch validators for TerraC: %w", err)
	}

	validators := validatorsResponse.Validators
	fmt.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		url := config.GetTerracConfig().API + "/cosmos/staking/v1beta1/validators/" + validator.OperatorAddress + "/delegations?pagination.limit=" + strconv.Itoa(config.LimitPerPage) + "&pagination.count_total=true"
		delegations, total, err := utils.FetchDelegations(url)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to fetch delegations for TerraC: %w", err)
		}
		fmt.Println(validator.OperatorAddress)
		fmt.Println("Response ", len(delegations))
		fmt.Println("Terrac validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegations...)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")

	apiURL := config.APICoingecko + config.GetTerracConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd, err := utils.FetchTokenPrice(apiURL, config.GetTerracConfig().CoinID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch TerraC token price: %w", err)
	}
	tokenIn20Usd := usd.QuoTruncate(tokenInUsd)

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenDelegate := sdkmath.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}
	eveAirdrop := sdkmath.LegacyMustNewDecFromStr(config.EveAirdrop)
	testAmount, err := sdkmath.LegacyNewDecFromStr("0")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to convert string to dec: %w", err)
	}
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetTerracConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address, err := utils.ConvertBech32Address(delegator.Delegation.DelegatorAddress)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetTerracConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println("Terrac balance: ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}
