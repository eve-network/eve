package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetLastValidators(sdk.Context) []staking.Validator
	HistoricalEntries(sdk.Context) uint32
	UnbondingTime(sdk.Context) time.Duration
}
