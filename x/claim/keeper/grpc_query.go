package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/eve-network/eve/x/claim/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) ClaimRecord(goCtx context.Context, req *types.QueryClaimRecordRequest) (*types.QueryClaimRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	claimRecord, found := k.GetClaimRecord(ctx, address)
	if !found {
		return nil, errors.Wrap(errors.ErrInvalidAddress, "Address not initial claim yet")
	}

	return &types.QueryClaimRecordResponse{
		ClaimRecord: claimRecord,
	}, nil

}

func (k Keeper) Claimable(goCtx context.Context, req *types.QueryTotalClaimableRequest) (*types.QueryTotalClaimableResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	claimable, found := k.GetClaimableInfo(ctx, address)
	if found {
		return &types.QueryTotalClaimableResponse{
			Amount: claimable.Amount,
		}, nil
	}

	claimRecord, found := k.GetClaimRecord(ctx, address)
	if !found {
		return nil, errors.Wrap(errors.ErrInvalidAddress, "Address not in airdrop list")
	}

	var amount sdk.Coins
	if !claimRecord.ClaimCompleted {
		amount = claimRecord.ClaimAble.Amount
	}

	return &types.QueryTotalClaimableResponse{
		Amount: amount,
	}, nil
}
