package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/notional-labs/eve/x/tokenfactory/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DenomAll(c context.Context, req *types.QueryAllDenomRequest) (*types.QueryAllDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var denoms []types.Denom
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	denomStore := prefix.NewStore(store, types.GetKeyPrefix(types.DenomKeyPrefix))

	pageRes, err := query.Paginate(denomStore, req.Pagination, func(key []byte, value []byte) error {
		var denom types.Denom
		if err := k.cdc.Unmarshal(value, &denom); err != nil {
			return err
		}

		denoms = append(denoms, denom)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDenomResponse{Denom: denoms, Pagination: pageRes}, nil
}

func (k Keeper) Denom(c context.Context, req *types.QueryGetDenomRequest) (*types.QueryGetDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetDenom(
		ctx,
		req.Denom,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetDenomResponse{Denom: val}, nil
}
