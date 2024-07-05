package chains

import (
	"context"
	"fmt"
	"log"

	"github.com/eve-network/eve/airdrop/config"
	"github.com/eve-network/eve/airdrop/utils"
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
		log.Printf("Failed to get latest height for Neutron: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to get latest height for Neutron: %w", err)
	}

	addresses, err := fetchBalance(blockHeight)
	if err != nil {
		log.Printf("Failed to get balance for Neutron: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch balance for Neutron: %w", err)
	}

	usd := sdkmath.LegacyMustNewDecFromStr("20")
	tokenInUsd, err := utils.FetchTokenPrice(config.GetNeutronConfig().CoinID)
	if err != nil {
		log.Printf("Failed to get Neutron token price: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to fetch Neutron token price: %w", err)
	}
	tokenIn20Usd := usd.Quo(tokenInUsd)

	rewardInfo := make([]config.Reward, 0, len(addresses))
	balanceInfo := make([]banktypes.Balance, 0, len(addresses))

	totalTokenBalance := sdkmath.NewInt(0)
	for _, address := range addresses {
		totalTokenBalance = totalTokenBalance.Add(address.Balance.Amount)
	}

	eveAirdrop, err := sdkmath.LegacyNewDecFromStr(config.EveAirdrop)
	if err != nil {
		log.Printf("Failed to convert EveAirdrop string to dec: %v", err)
		return nil, nil, 0, fmt.Errorf("failed to convert EveAirdrop string to dec: %w", err)
	}

	for _, address := range addresses {
		if sdkmath.LegacyNewDecFromInt(address.Balance.Amount).LT(tokenIn20Usd) {
			continue
		}

		eveAirdropAmount := eveAirdrop.MulInt64(int64(config.GetNeutronConfig().Percent)).QuoInt64(100).MulInt(address.Balance.Amount).QuoInt(totalTokenBalance)
		eveBech32Address, err := utils.ConvertBech32Address(address.Address)
		if err != nil {
			log.Printf("Failed to convert Neutron bech32 address: %v", err)
			return nil, nil, 0, fmt.Errorf("failed to convert Bech32Address: %w", err)
		}

		reward := config.Reward{
			Address:         address.Address,
			EveAddress:      eveBech32Address,
			Token:           address.Balance.Amount.ToLegacyDec(),
			EveAirdropToken: eveAirdropAmount,
			ChainID:         config.GetNeutronConfig().ChainID,
		}
		rewardInfo = append(rewardInfo, reward)

		balance := banktypes.Balance{
			Address: eveBech32Address,
			Coins:   sdk.NewCoins(sdk.NewCoin("eve", eveAirdropAmount.TruncateInt())),
		}
		balanceInfo = append(balanceInfo, balance)
	}

	return balanceInfo, rewardInfo, len(balanceInfo), nil
}

func fetchBalance(blockHeight string) ([]*banktypes.DenomOwner, error) {
	grpcAddr := config.GetNeutronConfig().GRPCAddr
	grpcConn, err := utils.SetupGRPCConnection(grpcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC Neutron: %w", err)
	}
	defer grpcConn.Close()

	bankClient := banktypes.NewQueryClient(grpcConn)
	addressInfo := []*banktypes.DenomOwner{}
	var paginationKey []byte
	step := 5000

	for {
		log.Printf("Fetching addresses %d to %d", len(addressInfo), len(addressInfo)+step)

		addresses, err := bankClient.DenomOwners(
			metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, blockHeight),
			&banktypes.QueryDenomOwnersRequest{
				Denom: "untrn",
				Pagination: &query.PageRequest{
					Limit:      uint64(step),
					Key:        paginationKey,
					CountTotal: true,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query all Neutron's denom owners: %w", err)
		}

		addressInfo = append(addressInfo, addresses.DenomOwners...)
		paginationKey = addresses.Pagination.NextKey
		if len(paginationKey) == 0 {
			break
		}
	}

	return addressInfo, nil
}
