package wasmbinding

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/eve/app"
	"github.com/notional-labs/eve/wasmbinding/bindings"
	"github.com/notional-labs/eve/x/gamm/pool-models/balancer"
)

// we must pay this many uosmo for every pool we create
var poolFee int64 = 1000000000

var defaultFunds = sdk.NewCoins(
	sdk.NewInt64Coin("uatom", 333000000),
	sdk.NewInt64Coin("uosmo", 555000000+2*poolFee),
	sdk.NewInt64Coin("ustar", 999000000),
)

func SetupCustomApp(t *testing.T, addr sdk.AccAddress) (*app.EveApp, sdk.Context) {
	osmosis, ctx := CreateTestInput()
	wasmKeeper := osmosis.WasmKeeper

	storeReflectCode(t, ctx, osmosis, addr)

	cInfo := wasmKeeper.GetCodeInfo(ctx, 1)
	require.NotNil(t, cInfo)

	return osmosis, ctx
}

func TestQueryFullDenom(t *testing.T) {
	actor := RandomAccountAddress()
	osmosis, ctx := SetupCustomApp(t, actor)

	reflect := instantiateReflectContract(t, ctx, osmosis, actor)
	require.NotEmpty(t, reflect)

	// query full denom
	query := bindings.OsmosisQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "ustart",
		},
	}
	resp := bindings.FullDenomResponse{}
	queryCustom(t, ctx, osmosis, reflect, query, &resp)

	expected := fmt.Sprintf("factory/%s/ustart", reflect.String())
	require.EqualValues(t, expected, resp.Denom)
}

func fundAccount(t *testing.T, ctx sdk.Context, osmosis *app.EveApp, addr sdk.AccAddress, coins sdk.Coins) {
	err := simapp.FundAccount(
		osmosis.BankKeeper,
		ctx,
		addr,
		coins,
	)
	require.NoError(t, err)
}

func preparePool(t *testing.T, ctx sdk.Context, osmosis *app.EveApp, addr sdk.AccAddress, funds []sdk.Coin) uint64 {
	var assets []balancer.PoolAsset
	for _, coin := range funds {
		assets = append(assets, balancer.PoolAsset{
			Weight: sdk.NewInt(100),
			Token:  coin,
		})
	}

	poolParams := balancer.PoolParams{
		SwapFee: sdk.NewDec(0),
		ExitFee: sdk.NewDec(0),
	}

	msg := balancer.NewMsgCreateBalancerPool(addr, poolParams, assets, "")
	poolId, err := osmosis.GAMMKeeper.CreatePool(ctx, &msg)
	require.NoError(t, err)
	return poolId
}
