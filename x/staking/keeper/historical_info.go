package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	dbm "github.com/tendermint/tm-db"

	// apptypes "github.com/notional-labs/eve/types"
	"github.com/notional-labs/eve/app/helpers"
	staking "github.com/notional-labs/eve/x/staking/types"
)

// Forked from Cosmos SDK x/staking/keeper/historical_info.go v0.42.9
// Forking was required to prevent modifications to the App_State for every call to BeginBlock.

const StoreKey = "history"

var (
	historyKeyprefix = []byte("emstaking/hist")
)

type HistoryKeeper struct {
	StoreKey      sdk.Store
	cdc           codec.BinaryCodec
	stakingKeeper staking.StakingKeeper
	database      dbm.DB
}

func getHistoricalInfoKey(height int64) []byte {
	key := types.GetHistoricalInfoKey(height)
	key = append(historyKeyprefix, key...)

	return key
}

// NewKeeper creates a new staking Keeper instance
func NewHistoryKeeper(cdc codec.Codec, key sdk.Store, stakingkeeper staking.StakingKeeper, db dbm.DB) HistoryKeeper {
	return HistoryKeeper{
		StoreKey:      key,
		cdc:           cdc,
		stakingKeeper: stakingkeeper,
		database:      db,
	}
}

// UnbondingTime is necessary to fulfill the interface required by the IBC module, which is the primary user of the HistoricalInfo.
func (k HistoryKeeper) UnbondingTime(ctx sdk.Context) time.Duration {
	return k.stakingKeeper.UnbondingTime(ctx)
}

// GetHistoricalInfo gets the historical info at a given height
func (k HistoryKeeper) GetHistoricalInfo(ctx sdk.Context, height int64) (types.HistoricalInfo, bool) {
	key := getHistoricalInfoKey(height)

	value, _ := k.database.Get(key)
	if value == nil {
		return types.HistoricalInfo{}, false
	}

	return types.MustUnmarshalHistoricalInfo(k.cdc, value), true
}

// SetHistoricalInfo sets the historical info at a given height
func (k HistoryKeeper) SetHistoricalInfo(ctx sdk.Context, height int64, hi *types.HistoricalInfo) {
	batch := helpers.GetCurrentBatch(ctx)
	if batch == nil {
		panic("batch object not found")
	}

	key := getHistoricalInfoKey(height)
	value := k.cdc.MustMarshal(hi)
	batch.Set(key, value)
}

// DeleteHistoricalInfo deletes the historical info at a given height
func (k HistoryKeeper) DeleteHistoricalInfo(ctx sdk.Context, height int64) {
	batch := helpers.GetCurrentBatch(ctx)
	if batch == nil {
		panic("batch object not found")
	}

	key := getHistoricalInfoKey(height)
	batch.Delete(key)
}

// IterateHistoricalInfo provides an interator over all stored HistoricalInfo
//
//	objects. For each HistoricalInfo object, cb will be called. If the cb returns
//
// true, the iterator will close and stop.
func (k HistoryKeeper) IterateHistoricalInfo(ctx sdk.Context, cb func(types.HistoricalInfo) bool) {
	iterator, err := k.database.Iterator(historyKeyprefix, sdk.PrefixEndBytes(historyKeyprefix))
	if err != nil {
		return
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		histInfo := types.MustUnmarshalHistoricalInfo(k.cdc, iterator.Value())
		if cb(histInfo) {
			break
		}
	}
}

// GetAllHistoricalInfo returns all stored HistoricalInfo objects.
func (k HistoryKeeper) GetAllHistoricalInfo(ctx sdk.Context) []types.HistoricalInfo {
	var infos []types.HistoricalInfo

	k.IterateHistoricalInfo(ctx, func(histInfo types.HistoricalInfo) bool {
		infos = append(infos, histInfo)
		return false
	})

	return infos
}

// TrackHistoricalInfo saves the latest historical-info and deletes the oldest
// heights that are below pruning height
func (k HistoryKeeper) TrackHistoricalInfo(ctx sdk.Context) {
	entryNum := k.stakingKeeper.HistoricalEntries(ctx)

	// Prune store to ensure we only have parameter-defined historical entries.
	// In most cases, this will involve removing a single historical entry.
	// In the rare scenario when the historical entries gets reduced to a lower value k'
	// from the original value k. k - k' entries must be deleted from the store.
	// Since the entries to be deleted are always in a continuous range, we can iterate
	// over the historical entries starting from the most recent version to be pruned
	// and then return at the first empty entry.
	for i := ctx.BlockHeight() - int64(entryNum); i >= 0; i-- {
		_, found := k.GetHistoricalInfo(ctx, i)
		if found {
			k.DeleteHistoricalInfo(ctx, i)
		} else {
			break
		}
	}

	// if there is no need to persist historicalInfo, return
	if entryNum == 0 {
		return
	}

	// Create HistoricalInfo struct
	lastVals := k.stakingKeeper.GetLastValidators(ctx)
	historicalEntry := types.NewHistoricalInfo(ctx.BlockHeader(), lastVals, sdk.DefaultPowerReduction)

	// Set latest HistoricalInfo at current height
	k.SetHistoricalInfo(ctx, ctx.BlockHeight(), &historicalEntry)
}
