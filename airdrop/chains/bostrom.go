package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Bostrom() ([]banktypes.Balance, []config.Reward, int, error) {
	delegators := []stakingtypes.DelegationResponse{}

	rpc := config.GetBostromConfig().API + "/cosmos/staking/v1beta1/validators?pagination.limit=" + strconv.Itoa(config.LimitPerPage) + "&pagination.count_total=true"
	validatorsResponse, err := utils.FetchValidators(rpc)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch validators for Bostrom: %w", err)
	}
	validators := validatorsResponse.Validators
	fmt.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		url := config.GetBostromConfig().API + "/cosmos/staking/v1beta1/validators/" + validator.OperatorAddress + "/delegations?pagination.limit=" + strconv.Itoa(config.LimitPerPage) + "&pagination.count_total=true"
		delegations, total, err := utils.FetchDelegations(url)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to fetch delegations for Bostrom: %w", err)
		}
		fmt.Println(validator.OperatorAddress)
		fmt.Println("Response ", len(delegations))
		fmt.Println("Bostrom validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegations...)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")

	apiURL := config.APICoingecko + config.GetBostromConfig().CoinID + "&vs_currencies=usd"
	fetchTokenPrice := utils.FetchTokenPriceWithRetry(fetchBostromTokenPrice)
	tokenInUsd, err := fetchTokenPrice(apiURL)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch Bostrom token price: %w", err)
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
	testAmount, _ := sdkmath.LegacyNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := utils.FindValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetBostromConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
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
			ChainID:         config.GetBostromConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println("Bostrom balance: ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}

func fetchBostromTokenPrice(apiURL string) (sdkmath.LegacyDec, error) {
	// Make a GET request to the API
	response, err := utils.MakeGetRequest(apiURL)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error making GET request to fetch Bostrom token price: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error reading response body for Bostrom token price: %w", err)
	}

	var data config.BostromPrice

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("error unmarshalling JSON for Bostrom token price: %w", err)
	}
	rawPrice := strings.Split(data.Token.USD.String(), "e-")
	base := rawPrice[0]
	power := rawPrice[1]
	powerInt, _ := strconv.ParseUint(power, 10, 64)
	baseDec, _ := sdkmath.LegacyNewDecFromStr(base)
	tenDec, _ := sdkmath.LegacyNewDecFromStr("10")
	tokenInUsd := baseDec.Quo(tenDec.Power(powerInt))
	return tokenInUsd, nil
}
