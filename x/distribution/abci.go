// This software is Copyright (c) 2019-2020 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package distribution

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/notional-labs/eve/app/helpers"

	// apptypes "github.com/notional-labs/eve/types"
	abci "github.com/tendermint/tendermint/abci/types"
	db "github.com/tendermint/tm-db"
)

var (
	previousProposerKey = []byte("emdistr/previousproposer")
)

const ModuleName = distrtypes.ModuleName

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

type DistributionKeeper interface {
	AllocateTokens(
		ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
		previousProposer sdk.ConsAddress, previousVotes []abci.VoteInfo,
	)
}

// Adapted from cosmos-sdk/x/distribution/abci.go
// A custom version was needed to keep the address of the previousProposer out of the consensus-state.

// set the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k DistributionKeeper, ak AccountKeeper, bk bankkeeper.ViewKeeper, db db.DB) {
	defer telemetry.ModuleMeasureSince(ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	batch := helpers.GetCurrentBatch(ctx)
	if batch == nil {
		panic("batch object not found")
	}

	// determine the total power signing the block
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousProposer, err := db.Get(previousProposerKey)
		if err != nil {
			panic(err)
		}

		feeCollector := ak.GetModuleAddress(auth.FeeCollectorName)
		coins := bk.GetAllBalances(ctx, feeCollector)

		// Only call AllocateTokens if there are in fact tokens to allocate.
		if !coins.IsZero() {
			k.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
		}
	}

	batch.Set(previousProposerKey, req.Header.ProposerAddress)
}
