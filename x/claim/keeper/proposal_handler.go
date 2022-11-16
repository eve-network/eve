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

}
