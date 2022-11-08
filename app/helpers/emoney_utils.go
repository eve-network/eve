package helpers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	db "github.com/tendermint/tm-db"
)

type contextKey uint8

const (
	_            contextKey = iota
	currentBatch contextKey = iota
)

func GetCurrentBatch(ctx sdk.Context) db.Batch {
	v, _ := ctx.Value(currentBatch).(db.Batch)
	return v
}

func WithCurrentBatch(ctx sdk.Context, batch db.Batch) sdk.Context {
	return ctx.WithValue(currentBatch, batch)
}
