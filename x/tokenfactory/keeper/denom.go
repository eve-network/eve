package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/eve/x/tokenfactory/types"
)

// SetDenom set a specific denom in the store from its index
func (k Keeper) SetDenom(ctx sdk.Context, denom types.Denom) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetKeyPrefix(types.DenomKeyPrefix))
	b := k.cdc.MustMarshal(&denom)
	store.Set(types.DenomKey(
		denom.Denom,
	), b)
}

// GetDenom returns a denom from its index
func (k Keeper) GetDenom(
	ctx sdk.Context,
	denom string,

) (val types.Denom, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetKeyPrefix(types.DenomKeyPrefix))

	b := store.Get(types.DenomKey(
		denom,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllDenom returns all denom
func (k Keeper) GetAllDenoms(ctx sdk.Context) (list []types.Denom) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetKeyPrefix(types.DenomKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Denom
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
