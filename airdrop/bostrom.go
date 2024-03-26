package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/eve-network/eve/airdrop/config"
)

func bostrom() ([]banktypes.Balance, []config.Reward) {
	delegators := []stakingtypes.DelegationResponse{}

	rpc := config.GetBostromConfig().API + "/cosmos/staking/v1beta1/validators?pagination.limit=" + strconv.Itoa(LIMIT_PER_PAGE) + "&pagination.count_total=true"
	validatorsResponse := fetchValidators(rpc)
	validators := validatorsResponse.Validators
	fmt.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		url := config.GetBostromConfig().API + "/cosmos/staking/v1beta1/validators/" + validator.OperatorAddress + "/delegations?pagination.limit=" + strconv.Itoa(LIMIT_PER_PAGE) + "&pagination.count_total=true"
		delegations, total := fetchDelegations(url)
		fmt.Println(validator.OperatorAddress)
		fmt.Println("Response ", len(delegations))
		fmt.Println("Validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegations...)
	}

	usd := math.LegacyMustNewDecFromStr("20")

	apiUrl := "https://api.coingecko.com/api/v3/simple/price?ids=" + config.GetBostromConfig().CoinId + "&vs_currencies=usd"
	tokenInUsd := fetchBostromTokenPrice(apiUrl)
	tokenIn20Usd := usd.QuoTruncate(tokenInUsd)

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenDelegate := math.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}
	eveAirdrop := math.LegacyMustNewDecFromStr(EVE_AIRDROP)
	testAmount, _ := math.LegacyNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfoCustomType(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetBostromConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address := convertBech32Address(delegator.Delegation.DelegatorAddress)
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdrop,
			ChainId:         config.GetBostromConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println("bostrom ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo
}

func fetchBostromTokenPrice(apiUrl string) math.LegacyDec {
	// Make a GET request to the API
	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		panic("")
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		panic("")
	}

	var data config.BostromPrice

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}
	rawPrice := strings.Split(data.Token.USD.String(), "e-")
	base := rawPrice[0]
	power := rawPrice[1]
	powerInt, _ := strconv.ParseUint(power, 10, 64)
	baseDec, _ := math.LegacyNewDecFromStr(base)
	tenDec, _ := math.LegacyNewDecFromStr("10")
	tokenInUsd := baseDec.Quo(tenDec.Power(powerInt))
	return tokenInUsd
}
