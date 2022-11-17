package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/eve-network/eve/x/claim/types"
)

func NewClaimProposalHandler(k Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.AirdropProposal:
			return handleAirdropProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized claim proposal content type: %T", c)
		}
	}
}

func handleAirdropProposal(ctx sdk.Context, k Keeper, p *types.AirdropProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	// Set params
	params := types.NewParams(true, p.ClaimDenom, p.AirdropStartTime, p.DurationUntilDecay, p.DurationOfDecay)
	err := k.SetParams(ctx, params)
	if err != nil {
		return err
	}
	// Set airdrop list
	err = k.appendClaimableList(ctx, p.AirdropList)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovInitialAirdrop,
	))

	return nil
}
