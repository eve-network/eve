package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// got to export genesis state from neutron and bostrom chain

const (
	EVE_AIRDROP    = "1000000000" // 1,000,000,000
	LIMIT_PER_PAGE = 100000000
)

func getValidators(stakingClient stakingtypes.QueryClient, block_height string) []stakingtypes.Validator {
	// Get validator
	var header metadata.MD
	var totalValidatorsResponse *stakingtypes.QueryValidatorsResponse
	totalValidatorsResponse, err := stakingClient.Validators(
		metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, block_height), // Add metadata to request
		&stakingtypes.QueryValidatorsRequest{
			Pagination: &query.PageRequest{
				Limit: LIMIT_PER_PAGE,
			},
		},
		grpc.Header(&header),
	)
	fmt.Println(err)
	validatorsInfo := totalValidatorsResponse.Validators
	return validatorsInfo
}

func main() {
	apiUrl := "https://api.coingecko.com/api/v3/simple/price?ids=" + config.GetBostromConfig().CoinId + "&vs_currencies=usd"
	fetchBostromTokenPrice(apiUrl)
	return
	balanceComposableInfo, rewardComposableInfo := bostrom()

	airdropMap := make(map[string]int)
	for _, info := range balanceComposableInfo {
		amount := airdropMap[info.Address]
		airdropMap[info.Address] = amount + int(info.Coins.AmountOf("eve").Int64())
	}

	balanceInfo := []banktypes.Balance{}
	checkBalance := 0
	for address, amount := range airdropMap {
		checkBalance += amount
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", math.NewInt(int64(amount)))),
		})
	}

	fmt.Println("Check balance: ", checkBalance)

	// Write delegations to file
	fileForDebug, _ := json.MarshalIndent(rewardComposableInfo, "", " ")
	_ = os.WriteFile("rewards.json", fileForDebug, 0644)

	fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	_ = os.WriteFile("balance.json", fileBalance, 0644)
}

func findValidatorInfo(validators []stakingtypes.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

func getLatestHeight(apiUrl string) string {
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

	// Print the response body
	var data config.NodeResponse

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}
	fmt.Println("Block height: ", data.Result.SyncInfo.LatestBlockHeight)
	return data.Result.SyncInfo.LatestBlockHeight
}

func convertBech32Address(otherChainAddress string) string {
	_, bz, _ := bech32.DecodeAndConvert(otherChainAddress)
	newBech32DelAddr, _ := bech32.ConvertAndEncode("eve", bz)
	return newBech32DelAddr
}
