package wasmbinding

import (
	"fmt"
	"testing"

	"github.com/notional-labs/eve/x/tokenfactory/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/eve/wasmbinding/bindings"
)

func TestCreateDenomMsg(t *testing.T) {
	creator := RandomAccountAddress()
	osmosis, ctx := SetupCustomApp(t, creator)

	lucky := RandomAccountAddress()
	reflect := instantiateReflectContract(t, ctx, osmosis, lucky)
	require.NotEmpty(t, reflect)

	// Fund reflect contract with 100 base denom creation fees
	reflectAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, osmosis, reflect, reflectAmount)

	msg := bindings.OsmosisMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "SUN",
	}}
	err := executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	// query the denom and see if it matches
	query := bindings.OsmosisQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp := bindings.FullDenomResponse{}
	queryCustom(t, ctx, osmosis, reflect, query, &resp)

	require.Equal(t, resp.Denom, fmt.Sprintf("factory/%s/SUN", reflect.String()))
}

func TestMintMsg(t *testing.T) {
	creator := RandomAccountAddress()
	osmosis, ctx := SetupCustomApp(t, creator)

	lucky := RandomAccountAddress()
	reflect := instantiateReflectContract(t, ctx, osmosis, lucky)
	require.NotEmpty(t, reflect)

	// Fund reflect contract with 100 base denom creation fees
	reflectAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, osmosis, reflect, reflectAmount)

	// lucky was broke
	balances := osmosis.BankKeeper.GetAllBalances(ctx, lucky)
	require.Empty(t, balances)

	// Create denom for minting
	msg := bindings.OsmosisMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "SUN",
	}}
	err := executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
	sunDenom := fmt.Sprintf("factory/%s/%s", reflect.String(), msg.CreateDenom.Subdenom)

	amount, ok := sdk.NewIntFromString("808010808")
	require.True(t, ok)
	msg = bindings.OsmosisMsg{MintTokens: &bindings.MintTokens{
		Denom:         sunDenom,
		Amount:        amount,
		MintToAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	balances = osmosis.BankKeeper.GetAllBalances(ctx, lucky)
	require.Len(t, balances, 1)
	coin := balances[0]
	require.Equal(t, amount, coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query := bindings.OsmosisQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp := bindings.FullDenomResponse{}
	queryCustom(t, ctx, osmosis, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)

	// mint the same denom again
	err = executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	balances = osmosis.BankKeeper.GetAllBalances(ctx, lucky)
	require.Len(t, balances, 1)
	coin = balances[0]
	require.Equal(t, amount.MulRaw(2), coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query = bindings.OsmosisQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp = bindings.FullDenomResponse{}
	queryCustom(t, ctx, osmosis, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)

	// now mint another amount / denom
	// create it first
	msg = bindings.OsmosisMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "MOON",
	}}
	err = executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
	moonDenom := fmt.Sprintf("factory/%s/%s", reflect.String(), msg.CreateDenom.Subdenom)

	amount = amount.SubRaw(1)
	msg = bindings.OsmosisMsg{MintTokens: &bindings.MintTokens{
		Denom:         moonDenom,
		Amount:        amount,
		MintToAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	balances = osmosis.BankKeeper.GetAllBalances(ctx, lucky)
	require.Len(t, balances, 2)
	coin = balances[0]
	require.Equal(t, amount, coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query = bindings.OsmosisQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "MOON",
		},
	}
	resp = bindings.FullDenomResponse{}
	queryCustom(t, ctx, osmosis, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)

	// and check the first denom is unchanged
	coin = balances[1]
	require.Equal(t, amount.AddRaw(1).MulRaw(2), coin.Amount)
	require.Contains(t, coin.Denom, "factory/")

	// query the denom and see if it matches
	query = bindings.OsmosisQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "SUN",
		},
	}
	resp = bindings.FullDenomResponse{}
	queryCustom(t, ctx, osmosis, reflect, query, &resp)

	require.Equal(t, resp.Denom, coin.Denom)
}

func TestBurnMsg(t *testing.T) {
	creator := RandomAccountAddress()
	osmosis, ctx := SetupCustomApp(t, creator)

	lucky := RandomAccountAddress()
	reflect := instantiateReflectContract(t, ctx, osmosis, lucky)
	require.NotEmpty(t, reflect)

	// Fund reflect contract with 100 base denom creation fees
	reflectAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)))
	fundAccount(t, ctx, osmosis, reflect, reflectAmount)

	// lucky was broke
	balances := osmosis.BankKeeper.GetAllBalances(ctx, lucky)
	require.Empty(t, balances)

	// Create denom for minting
	msg := bindings.OsmosisMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "SUN",
	}}
	err := executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
	sunDenom := fmt.Sprintf("factory/%s/%s", reflect.String(), msg.CreateDenom.Subdenom)

	amount, ok := sdk.NewIntFromString("808010808")
	require.True(t, ok)

	msg = bindings.OsmosisMsg{MintTokens: &bindings.MintTokens{
		Denom:         sunDenom,
		Amount:        amount,
		MintToAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)

	// can't burn from different address
	msg = bindings.OsmosisMsg{BurnTokens: &bindings.BurnTokens{
		Denom:           sunDenom,
		Amount:          amount,
		BurnFromAddress: lucky.String(),
	}}
	err = executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.Error(t, err)

	// lucky needs to send balance to reflect contract to burn it
	luckyBalance := osmosis.BankKeeper.GetAllBalances(ctx, lucky)
	err = osmosis.BankKeeper.SendCoins(ctx, lucky, reflect, luckyBalance)
	require.NoError(t, err)

	msg = bindings.OsmosisMsg{BurnTokens: &bindings.BurnTokens{
		Denom:           sunDenom,
		Amount:          amount,
		BurnFromAddress: reflect.String(),
	}}
	err = executeCustom(t, ctx, osmosis, reflect, lucky, msg, sdk.Coin{})
	require.NoError(t, err)
}
