package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/eve-network/eve/x/claim/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) InitialClaim(goCtx context.Context, msg *types.MsgInitialClaim) (*types.MsgInitialClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	params := k.GetParams(ctx)
	if !params.IsAirdropEnabled(ctx.BlockTime()) {
		return nil, types.ErrAirdropNotEnabled
	}
	claimInfo, found := k.GetClaimableInfo(ctx, sender)
	if !found {
		return nil, errors.Wrapf(errors.ErrInvalidAddress, "Invalid address: unclaimable or already initial claim %s", msg.Sender)
	}
	claimRecord := types.ClaimRecord{
		ClaimAble:      claimInfo,
		ClaimCompleted: false,
	}
	err = k.SetClaimRecord(ctx, claimRecord)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ClaimModuleEventType),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgInitialClaimResponse{
		ClaimableAmount: claimRecord.ClaimAble.Amount,
	}, nil
}

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	params := k.GetParams(ctx)
	if !params.IsAirdropEnabled(ctx.BlockTime()) {
		return nil, types.ErrAirdropNotEnabled
	}
	claimRecord, found := k.GetClaimRecord(ctx, sender)
	if !found {
		return nil, errors.Wrapf(errors.ErrInvalidAddress, "Invalid address: Not initialize %s", msg.Sender)
	}
	if claimRecord.ClaimCompleted {
		return nil, errors.Wrapf(errors.ErrInvalidAddress, "Invalid address: Already claim %s", msg.Sender)
	}
	coins, err := k.ClaimCoins(ctx, claimRecord.ClaimAble)
	if err != nil {
		return nil, err
	}
	claimRecord.ClaimCompleted = true
	err = k.SetClaimRecord(ctx, claimRecord)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ClaimModuleEventType),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgClaimResponse{
		Address:       msg.Sender,
		ClaimedAmount: coins,
	}, nil
}
