package keeper_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	// apptypes "github.com/notional-labs/eve/types"
	"github.com/notional-labs/eve/app/helpers"
	"github.com/notional-labs/eve/x/staking/keeper"
	emtypes "github.com/notional-labs/eve/x/staking/types"
)

func TestHistoricalInfo(t *testing.T) {
	ctx, batch, database, historyKeeper := setup(t)

	validatorCount := 50
	validators := make([]types.Validator, validatorCount)

	for i := 0; i < validatorCount; i++ {
		validators[i] = teststaking.NewValidator(t, sdk.ValAddress(fmt.Sprintf("validator%v", i)), secp256k1.GenPrivKey().PubKey())
	}

	hi := types.NewHistoricalInfo(ctx.BlockHeader(), validators, sdk.DefaultPowerReduction)
	historyKeeper.SetHistoricalInfo(ctx, 2, &hi)

	batch.WriteSync()
	batch = database.NewBatch()
	ctx = helpers.WithCurrentBatch(ctx, batch)

	recv, found := historyKeeper.GetHistoricalInfo(ctx, 2)
	require.True(t, found, "HistoricalInfo not found after set")
	require.Equal(t, hi, recv, "HistoricalInfo not equal")
	valset := validatorSet{
		Set: types.ValidatorsByVotingPower(recv.Valset),
	}
	require.True(t, sort.IsSorted(valset), "HistoricalInfo validators is not sorted")

	historyKeeper.DeleteHistoricalInfo(ctx, 2)

	err := batch.WriteSync()
	require.NoError(t, err)

	recv, found = historyKeeper.GetHistoricalInfo(ctx, 2)
	require.False(t, found, "HistoricalInfo found after delete")
	require.Equal(t, types.HistoricalInfo{}, recv, "HistoricalInfo is not empty")
}

func TestGetAllHistoricalInfo(t *testing.T) {
	ctx, batch, _, historyKeeper := setup(t)

	valSet := []types.Validator{
		teststaking.NewValidator(t, sdk.ValAddress("val0"), secp256k1.GenPrivKey().PubKey()),
		teststaking.NewValidator(t, sdk.ValAddress("val1"), secp256k1.GenPrivKey().PubKey()),
	}

	header1 := tmproto.Header{ChainID: "HelloChain", Height: 10}
	header2 := tmproto.Header{ChainID: "HelloChain", Height: 11}
	header3 := tmproto.Header{ChainID: "HelloChain", Height: 12}

	hist1 := types.HistoricalInfo{Header: header1, Valset: valSet}
	hist2 := types.HistoricalInfo{Header: header2, Valset: valSet}
	hist3 := types.HistoricalInfo{Header: header3, Valset: valSet}

	expHistInfos := []types.HistoricalInfo{hist1, hist2, hist3}

	for i, hi := range expHistInfos {
		historyKeeper.SetHistoricalInfo(ctx, int64(10+i), &hi)
	}

	batch.WriteSync()

	infos := historyKeeper.GetAllHistoricalInfo(ctx)
	require.NotEmpty(t, infos)
	require.Equal(t, expHistInfos, infos)
}

func setup(t *testing.T) (sdk.Context, db.Batch, db.DB, keeper.HistoryKeeper) {
	sdk.DefaultPowerReduction = sdk.OneInt()
	ms := store.NewCommitMultiStore(dbm.NewMemDB())
	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	ctx := sdk.NewContext(ms, tmproto.Header{ChainID: "HelloChain"}, true, log.NewNopLogger())

	database := db.NewMemDB()
	batch := database.NewBatch()
	ctx = helpers.WithCurrentBatch(ctx, batch)

	key := sdk.NewKVStoreKey(keeper.StoreKey)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	historyKeeper := keeper.NewHistoryKeeper(marshaler, ctx.KVStore(key), mockStakingKeeper{}, database)

	return ctx, batch, database, historyKeeper
}

var _ emtypes.StakingKeeper = mockStakingKeeper{}

type mockStakingKeeper struct{}

func (m mockStakingKeeper) GetLastValidators(ctx sdk.Context) (_ []types.Validator) {
	return
}

func (m mockStakingKeeper) HistoricalEntries(ctx sdk.Context) uint32 {
	return 1000
}

func (m mockStakingKeeper) UnbondingTime(ctx sdk.Context) time.Duration {
	return 24 * 21 * time.Hour
}

type validatorSet struct {
	Set types.ValidatorsByVotingPower
}

func (v validatorSet) Less(i, j int) bool {
	return v.Set.Less(i, j, sdk.DefaultPowerReduction)
}

func (v validatorSet) Len() int {
	return v.Set.Len()
}

func (v validatorSet) Swap(i, j int) {
	v.Set.Swap(i, j)
}
