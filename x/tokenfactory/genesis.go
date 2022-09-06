package tokenfactory

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/eve/x/tokenfactory/keeper"
	"github.com/notional-labs/eve/x/tokenfactory/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the denom
	for _, elem := range genState.DenomList {
		k.SetDenom(ctx, elem)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.DenomList = k.GetAllDenoms(ctx)

	return genesis
}
