package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/eve-network/eve/x/claim/types"
	"github.com/gogo/protobuf/proto"
)

// CreateModuleAccount creates module account of airdrop module
func (k Keeper) CreateModuleAccount(ctx sdk.Context, amount sdk.Coin) {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
}

func (k Keeper) ClaimCoins(ctx sdk.Context, claimable types.Claimable) (sdk.Coins, error) {

	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(claimable.Amount...))
	if err != nil {
		panic(err)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(claimable.Address), claimable.Amount)
	if err != nil {
		return nil, err
	}

	return claimable.Amount, nil
}

func (k Keeper) appendClaimableList(ctx sdk.Context, claimables []types.Claimable) error {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ClaimableListStorePrefix)

	for _, claimable := range claimables {
		accAddress, err := sdk.AccAddressFromBech32(claimable.Address)
		if err != nil {
			return err
		}

		bz, err := proto.Marshal(&claimable)
		if err != nil {
			return err
		}

		prefixStore.Set(accAddress, bz)
	}
	return nil
}

func (k Keeper) GetClaimableInfo(ctx sdk.Context, accAddress sdk.AccAddress) (types.Claimable, bool) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ClaimableListStorePrefix)

	var claimableInfo types.Claimable
	bz := prefixStore.Get(accAddress)
	if bz == nil {
		return claimableInfo, false
	}

	err := proto.Unmarshal(bz, &claimableInfo)
	if err != nil {
		panic(err)
	}

	prefixStore.Delete(accAddress)
	return claimableInfo, true
}

func (k Keeper) GetClaimRecord(ctx sdk.Context, address sdk.AccAddress) (types.ClaimRecord, bool) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ClaimRecordsStorePrefix)

	var claimRecord types.ClaimRecord
	bz := prefixStore.Get(address)
	if bz == nil {
		return claimRecord, false
	}

	err := proto.Unmarshal(bz, &claimRecord)
	if err != nil {
		panic(err)
	}

	return claimRecord, true
}

// SetClaimables set claimable amount from balances object
func (k Keeper) SetClaimRecords(ctx sdk.Context, claimRecords []types.ClaimRecord) error {
	for _, claimRecord := range claimRecords {
		err := k.SetClaimRecord(ctx, claimRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetClaimRecord sets a claim record for an address in store
func (k Keeper) SetClaimRecord(ctx sdk.Context, claimRecord types.ClaimRecord) error {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ClaimRecordsStorePrefix)

	bz, err := proto.Marshal(&claimRecord)
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(claimRecord.ClaimAble.Address)
	if err != nil {
		return err
	}

	prefixStore.Set(addr, bz)
	return nil
}

func (k Keeper) ClaimRecords(ctx sdk.Context) []types.ClaimRecord {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ClaimRecordsStorePrefix)

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	var claimRecords []types.ClaimRecord
	for ; iterator.Valid(); iterator.Next() {
		claimRecord := types.ClaimRecord{}

		err := proto.Unmarshal(iterator.Value(), &claimRecord)
		if err != nil {
			panic(err)
		}

		claimRecords = append(claimRecords, claimRecord)
	}

	return claimRecords
}

func (k Keeper) clearClaimRecords(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ClaimRecordsStorePrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

func (k Keeper) clearInitialClaimables(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ClaimableListStorePrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

func (k Keeper) EndAirdrop(ctx sdk.Context) error {
	k.clearInitialClaimables(ctx)
	k.clearInitialClaimables(ctx)
	return nil
}
