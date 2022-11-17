package claim

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eve-network/eve/x/claim/keeper"
	"github.com/eve-network/eve/x/claim/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	if !params.IsAirdropEnabled(ctx.BlockTime()) {
		return
	}
	// End Airdrop
	goneTime := ctx.BlockTime().Sub(params.AirdropStartTime)
	if goneTime > params.DurationUntilDecay+params.DurationOfDecay {
		// airdrop time passed
		err := k.EndAirdrop(ctx)
		if err != nil {
			panic(err)
		}
		// Clear params
		newParams := types.DefaultParams()
		err = k.SetParams(ctx, newParams)
		if err != nil {
			panic(err)
		}
	}
}
