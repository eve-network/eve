package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/notional-labs/eve/x/tokenfactory/types"
)

func (k msgServer) MintAndSendTokens(goCtx context.Context, msg *types.MsgMintAndSendTokens) (*types.MsgMintAndSendTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetDenom(
		ctx,
		msg.Denom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "denom does not exist")
	}

	// Checks if the the msg owner is the same as the current owner
	if msg.Owner != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if valFound.Supply+msg.Amount > valFound.MaxSupply {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Cannot mint more than Max Supply")
	}
	moduleAcct := k.accountKeeper.GetModuleAddress(types.ModuleName)

	recipientAddress, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	var mintCoins sdk.Coins

	mintCoins = mintCoins.Add(sdk.NewCoin(msg.Denom, sdk.NewInt(int64(msg.Amount))))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return nil, err
	}
	if err := k.bankKeeper.SendCoins(ctx, moduleAcct, recipientAddress, mintCoins); err != nil {
		return nil, err
	}

	var denom = types.Denom{
		Owner:              valFound.Owner,
		Name:               valFound.Name,
		Denom:              valFound.Denom,
		MaxSupply:          valFound.MaxSupply,
		Supply:             valFound.Supply + msg.Amount,
		Precision:          valFound.Precision,
		CanChangeMaxSupply: valFound.CanChangeMaxSupply,
	}

	k.SetDenom(
		ctx,
		denom,
	)
	return &types.MsgMintAndSendTokensResponse{}, nil
}
