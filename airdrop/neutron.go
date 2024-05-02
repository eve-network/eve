package main

// code = Unimplemented desc = unknown service cosmos.staking.v1beta1.Query
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eve-network/eve/airdrop/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func neutron() ([]banktypes.Balance, []config.Reward) {
	blockHeight := getLatestHeight(config.GetNeutronConfig().RPC + "/status")
	addresses, total := fetchBalance(blockHeight)
	fmt.Println("Response ", len(addresses))
	fmt.Println("Total ", total)

	usd, _ := math.LegacyNewDecFromStr("20")

	apiURL := APICoingecko + config.GetNeutronConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd := fetchNeutronTokenPrice(apiURL)
	tokenIn20Usd := usd.Quo(tokenInUsd)
	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenBalance, _ := math.NewIntFromString("0")
	for _, address := range addresses {
		if math.LegacyNewDecFromInt(address.Balance.Amount).LT(tokenIn20Usd) {
			continue
		}
		totalTokenBalance = totalTokenBalance.Add(address.Balance.Amount)
	}
	eveAirdrop := math.LegacyMustNewDecFromStr(EveAirdrop)
	testAmount, _ := math.LegacyNewDecFromStr("0")
	for _, address := range addresses {
		if math.LegacyNewDecFromInt(address.Balance.Amount).LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetNeutronConfig().Percent))).QuoInt64(100).MulInt(address.Balance.Amount).QuoInt(totalTokenBalance)
		eveBech32Address := convertBech32Address(address.Address)
		rewardInfo = append(rewardInfo, config.Reward{
			Address:         address.Address,
			EveAddress:      eveBech32Address,
			Token:           address.Balance.Amount.ToLegacyDec(),
			EveAirdropToken: eveAirdrop,
			ChainID:         config.GetNeutronConfig().ChainID,
		})
		testAmount = eveAirdrop.Add(testAmount)
		balanceInfo = append(balanceInfo, banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdrop.TruncateInt())),
		})
	}
	fmt.Println("Neutron ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo
}

func fetchBalance(blockHeight string) ([]*banktypes.DenomOwner, uint64) {
	grpcAddr := config.GetNeutronConfig().GRPCAddr
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())))
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()
	bankClient := banktypes.NewQueryClient(grpcConn)
	var header metadata.MD
	var addresses *banktypes.QueryDenomOwnersResponse // QueryValidatorDelegationsResponse
	var paginationKey []byte
	addressInfo := []*banktypes.DenomOwner{}
	step := 5000
	total := uint64(0)
	// Fetch addresses, 5000 at a time
	i := 0
	for {
		i++
		fmt.Println("Fetching addresses", step*i, "to", step*(i+1))
		addresses, err = bankClient.DenomOwners(
			metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight), // Add metadata to request
			&banktypes.QueryDenomOwnersRequest{
				Denom: "untrn",
				Pagination: &query.PageRequest{
					Limit:      uint64(step),
					Key:        paginationKey,
					CountTotal: true,
				},
			},
			grpc.Header(&header), // Retrieve header from response
		)
		fmt.Println("err: ", err)
		if total != 0 {
			total = addresses.Pagination.Total
		}
		addressInfo = append(addressInfo, addresses.DenomOwners...)
		paginationKey = addresses.Pagination.NextKey
		if len(paginationKey) == 0 {
			break
		}
	}
	return addressInfo, total
}

func fetchNeutronTokenPrice(apiURL string) math.LegacyDec {
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

	var data config.NeutronPrice

	// Unmarshal the JSON byte slice into the defined struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		panic("")
	}

	tokenInUsd := math.LegacyMustNewDecFromStr(data.Token.USD.String())
	return tokenInUsd
}
