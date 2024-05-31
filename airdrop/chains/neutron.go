package chains

// code = Unimplemented desc = unknown service cosmos.staking.v1beta1.Query
import (
	"context"
	"fmt"
	"log"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func Neutron() ([]banktypes.Balance, []config.Reward, int, error) {
	blockHeight, err := utils.GetLatestHeight(config.GetNeutronConfig().RPC + "/status")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Neutron: %w", err)
	}

	addresses, total, err := fetchBalance(blockHeight)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch balance for Neutron: %w", err)
	}
	log.Println("Response ", len(addresses))
	log.Println("Total ", total)

	usd, _ := sdkmath.LegacyNewDecFromStr("20")

	apiURL := config.APICoingecko + config.GetNeutronConfig().CoinID + "&vs_currencies=usd"
	tokenInUsd, err := utils.FetchTokenPrice(apiURL, config.GetNeutronConfig().CoinID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to fetch Neutron token price: %w", err)
	}
	tokenIn20Usd := usd.Quo(tokenInUsd)
	rewardInfo := []config.Reward{}
	balanceInfo := []banktypes.Balance{}

	totalTokenBalance, _ := sdkmath.NewIntFromString("0")
	for _, address := range addresses {
		totalTokenBalance = totalTokenBalance.Add(address.Balance.Amount)
	}
	eveAirdrop := sdkmath.LegacyMustNewDecFromStr(config.EveAirdrop)
	testAmount, _ := sdkmath.LegacyNewDecFromStr("0")
	for _, address := range addresses {
		if sdkmath.LegacyNewDecFromInt(address.Balance.Amount).LT(tokenIn20Usd) {
			continue
		}
		eveAirdrop := (eveAirdrop.MulInt64(int64(config.GetNeutronConfig().Percent))).QuoInt64(100).MulInt(address.Balance.Amount).QuoInt(totalTokenBalance)
		eveBech32Address, err := utils.ConvertBech32Address(address.Address)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}
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
	log.Println("Neutron balance: ", testAmount)
	// Write delegations to file
	// fileForDebug, _ := json.MarshalIndent(rewardInfo, "", " ")
	// _ = os.WriteFile("rewards.json", fileForDebug, 0644)

	// fileBalance, _ := json.MarshalIndent(balanceInfo, "", " ")
	// _ = os.WriteFile("balance.json", fileBalance, 0644)
	return balanceInfo, rewardInfo, len(balanceInfo), nil
}

func fetchBalance(blockHeight string) ([]*banktypes.DenomOwner, uint64, error) {
	grpcAddr := config.GetNeutronConfig().GRPCAddr
	grpcConn, err := utils.SetupGRPCConnection(grpcAddr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to connect to gRPC Neutron: %w", err)
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
		log.Println("Fetching addresses", step*i, "to", step*(i+1))
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
		if err != nil {
			return nil, 0, fmt.Errorf("failed to query all Neutron's denom owners: %w", err)
		}
		if total != 0 {
			total = addresses.Pagination.Total
		}
		addressInfo = append(addressInfo, addresses.DenomOwners...)
		paginationKey = addresses.Pagination.NextKey
		if len(paginationKey) == 0 {
			break
		}
	}
	return addressInfo, total, nil
}
