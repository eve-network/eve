package tokenfactory_test

import (
	"testing"

	keepertest "github.com/eve-network/eve/testutil/keeper"
	"github.com/eve-network/eve/x/tokenfactory"
	"github.com/eve-network/eve/x/tokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		DenomList: []types.Denom{
			{
				Denom: "0",
			},
			{
				Denom: "1",
			},
		},
	}

	k, ctx := keepertest.TokenfactoryKeeper(t)
	tokenfactory.InitGenesis(ctx, *k, genesisState)
	got := tokenfactory.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.DenomList, got.DenomList)
}
