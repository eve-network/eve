package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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
	balanceComposableInfo, rewardComposableInfo := terra()
	// Write delegations to file
	fileForDebug, _ := json.MarshalIndent(rewardComposableInfo, "", " ")
	_ = os.WriteFile("rewards.json", fileForDebug, 0644)

	fileBalance, _ := json.MarshalIndent(balanceComposableInfo, "", " ")
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
