package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// got to export genesis state from neutron and bostrom chain

const (
	EveAirdrop   = "1000000000" // 1,000,000,000
	LimitPerPage = 100000000
	Badkids      = "stars19jq6mj84cnt9p7sagjxqf8hxtczwc8wlpuwe4sh62w45aheseues57n420"
	Cryptonium   = "stars1g2ptrqnky5pu70r3g584zpk76cwqplyc63e8apwayau6l3jr8c0sp9q45u"
	APICoingecko = "https://api.coingecko.com/api/v3/simple/price?ids="
)

func getValidators(stakingClient stakingtypes.QueryClient, blockHeight string) []stakingtypes.Validator {
	// Get validator
	var header metadata.MD
	var totalValidatorsResponse *stakingtypes.QueryValidatorsResponse
	totalValidatorsResponse, err := stakingClient.Validators(
		metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight), // Add metadata to request
		&stakingtypes.QueryValidatorsRequest{
			Pagination: &query.PageRequest{
				Limit: LimitPerPage,
			},
		},
		grpc.Header(&header),
	)
	fmt.Println(err)
	validatorsInfo := totalValidatorsResponse.Validators
	return validatorsInfo
}

func main() {
	balanceAkashInfo, _ := akash()
	akashLength := len(balanceAkashInfo)

	balanceBostromInfo, _ := bostrom()
	bostromLength := len(balanceBostromInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceBostromInfo...)

	balanceCelestiaInfo, _ := celestia()
	celestiaLength := len(balanceCelestiaInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceCelestiaInfo...)

	balanceComposableInfo, _ := composable()
	composableLength := len(balanceComposableInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceComposableInfo...)

	balanceCosmosInfo, _ := cosmos()
	cosmosLength := len(balanceCosmosInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceCosmosInfo...)

	balanceNeutronInfo, _ := neutron()
	neutronLength := len(balanceNeutronInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceNeutronInfo...)

	balanceSentinelInfo, _ := sentinel()
	sentinelLength := len(balanceSentinelInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceSentinelInfo...)

	balanceStargazeInfo, _ := stargaze()
	stargazeLength := len(balanceStargazeInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceStargazeInfo...)

	balanceTerraInfo, _ := terra()
	terraLength := len(balanceTerraInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceTerraInfo...)

	balanceTerracInfo, _ := terrac()
	terracLength := len(balanceTerracInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceTerracInfo...)

	balanceBadKidsInfo, _ := cosmosnft(Badkids, int64(config.GetBadKidsConfig().Percent))
	badkidsLength := len(balanceBadKidsInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceBadKidsInfo...)

	balanceCryptoniumInfo, _ := cosmosnft(Cryptonium, int64(config.GetCryptoniumConfig().Percent))
	cryptoniumLength := len(balanceCryptoniumInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceCryptoniumInfo...)

	// need set coin type on Eve
	balanceMiladyInfo, _ := ethereumnft()
	miladyLength := len(balanceMiladyInfo)
	balanceAkashInfo = append(balanceAkashInfo, balanceMiladyInfo...)

	total := akashLength + bostromLength + celestiaLength + composableLength + cosmosLength + neutronLength + sentinelLength + stargazeLength + terraLength + terracLength + badkidsLength + cryptoniumLength + miladyLength
	fmt.Println("total: ", total)
	fmt.Println(len(balanceAkashInfo))

	airdropMap := make(map[string]int)
	for _, info := range balanceAkashInfo {
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

	// // Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardComposableInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	_ = os.WriteFile("balance.json", fileBalance, 0o600)
}

func findValidatorInfo(validators []stakingtypes.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

func getLatestHeight(apiURL string) string {
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

func fetchValidators(rpcURL string) config.ValidatorResponse {
	// Make a GET request to the API
	response, err := http.Get(rpcURL) //nolint
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

	var data config.ValidatorResponse

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}
	fmt.Println(data.Pagination.Total)
	return data
}

func findValidatorInfoCustomType(validators []config.Validator, address string) int {
	for key, v := range validators {
		if v.OperatorAddress == address {
			return key
		}
	}
	return -1
}

func fetchDelegations(rpcURL string) (stakingtypes.DelegationResponses, uint64) {
	// Make a GET request to the API
	response, err := http.Get(rpcURL) //nolint
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

	var data config.QueryValidatorDelegationsResponse

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}

	fmt.Println(data.Pagination.Total)
	total, _ := strconv.ParseUint(data.Pagination.Total, 10, 64)
	return data.DelegationResponses, total
}
