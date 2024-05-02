package main

// error max size response
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func celestia() ([]banktypes.Balance, []config.Reward, int) {
	blockHeight := getLatestHeight(config.GetCelestiaConfig().RPC + "/status")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env:", err)
		panic("")
	}
	grpcAddr := config.GetCelestiaConfig().GRPCAddr
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	delegators := []stakingtypes.DelegationResponse{}

	validators := getValidators(stakingClient, blockHeight)
	fmt.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		url := config.GetCelestiaConfig().API + "/cosmos/staking/v1beta1/validators/" + validator.OperatorAddress + "/delegations?pagination.limit=" + strconv.Itoa(LimitPerPage) + "&pagination.count_total=true"
		delegations, total := fetchDelegations(url)
		fmt.Println(validator.OperatorAddress)
		fmt.Println("Response ", len(delegations))
		fmt.Println("Celestia validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegations...)
	}

	usd := math.LegacyMustNewDecFromStr("20")

	apiURL := APICoingecko + config.GetCelestiaConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd := fetchCelestiaTokenPrice(apiURL)
	tokenIn20Usd := usd.QuoTruncate(tokenInUsd)

	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenDelegate := math.LegacyMustNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		totalTokenDelegate = totalTokenDelegate.Add(token)
	}
	eveAirdrop := math.LegacyMustNewDecFromStr(EveAirdrop)
	testAmount, _ := math.LegacyNewDecFromStr("0")
	for _, delegator := range delegators {
		validatorIndex := findValidatorInfo(validators, delegator.Delegation.ValidatorAddress)
		validatorInfo := validators[validatorIndex]
		token := (delegator.Delegation.Shares.MulInt(validatorInfo.Tokens)).QuoTruncate(validatorInfo.DelegatorShares)
		if token.LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetCelestiaConfig().Percent))).QuoInt64(100).Mul(token).QuoTruncate(totalTokenDelegate)
		eveBech32Address := convertBech32Address(delegator.Delegation.DelegatorAddress)
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         delegator.Delegation.DelegatorAddress,
			EveAddress:      eveBech32Address,
			Shares:          delegator.Delegation.Shares,
			Token:           token,
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetCelestiaConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println("Celestia ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo)
}

func fetchCelestiaTokenPrice(apiURL string) math.LegacyDec {
	// Make a GET request to the API
	response, err := http.Get(apiURL) //nolint
	if err != nil {
		fmt.Println("Error making GET request:", err)
		panic("")
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		panic("")
	}

	var data config.CelestiaPrice

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}

	tokenInUsd := math.LegacyMustNewDecFromStr(data.Token.USD.String())
	return tokenInUsd
}
