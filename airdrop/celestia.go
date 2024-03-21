package main

// error max size response
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/eve-network/eve/airdrop/config"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func celestia() ([]banktypes.Balance, []config.Reward) {
	block_height := getLatestHeight(config.GetCelestiaConfig().RPC + "/status")
	godotenv.Load()
	grpcAddr := config.GetCelestiaConfig().GRPCAddr
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())))
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	delegators := []stakingtypes.DelegationResponse{}

	validators := getValidators(stakingClient, block_height)
	fmt.Println("Validators: ", len(validators))
	for validatorIndex, validator := range validators {
		// var header metadata.MD
		// delegationsResponse, err := stakingClient.ValidatorDelegations(
		// 	metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, block_height), // Add metadata to request
		// 	&stakingtypes.QueryValidatorDelegationsRequest{
		// 		ValidatorAddr: validator.OperatorAddress,
		// 		Pagination: &query.PageRequest{
		// 			CountTotal: true,
		// 			Limit:      LIMIT_PER_PAGE,
		// 		},
		// 	},
		// 	grpc.Header(&header), // Retrieve header from response
		// )
		// stakingtypes.QueryValidatorDelegationsResponse
		rpcUrl := config.GetCelestiaConfig().API + "/validators/" + validator.String() + "/delegations?pagination.limit=" + string(LIMIT_PER_PAGE) + "&pagination.count_total=true"
		fmt.Println(rpcUrl)
		delegationsResponse := fetchDelegations(rpcUrl)
		fmt.Println("err: ", err)
		total := delegationsResponse.Pagination.Total
		fmt.Println("Response ", len(delegationsResponse.DelegationResponses))
		fmt.Println("Validator "+strconv.Itoa(validatorIndex)+" ", total)
		delegators = append(delegators, delegationsResponse.DelegationResponses...)
	}

	usd := math.LegacyMustNewDecFromStr("20")

	apiUrl := "https://api.coingecko.com/api/v3/simple/price?ids=" + config.GetCelestiaConfig().CoinId + "&vs_currencies=usd"
	tokenInUsd := fetchCelestiaTokenPrice(apiUrl)
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
	eveAirdrop := math.LegacyMustNewDecFromStr(EVE_AIRDROP)
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
			ChainId:         config.GetCelestiaConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println(testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo
}

func fetchDelegations(rpcUrl string) stakingtypes.QueryValidatorDelegationsResponse {
	// Make a GET request to the API
	response, err := http.Get(rpcUrl)
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

	var data stakingtypes.QueryValidatorDelegationsResponse
	fmt.Println(responseBody)
	// Unmarshal the JSON byte slice into the defined struct
	err = data.Unmarshal(responseBody)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}

	fmt.Println(data)
	return data
}

func fetchCelestiaTokenPrice(apiUrl string) math.LegacyDec {
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
