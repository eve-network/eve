package keeper

import (
	"context"

	"github.com/eve-network/eve/x/claim/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ModuleAccountBalance(context.Context, *types.QueryModuleAccountBalanceRequest) (*types.QueryModuleAccountBalanceResponse, error) {

}

func (k Keeper) Params(context.Context, *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {

}

func (k Keeper) ClaimRecord(context.Context, *types.QueryClaimRecordRequest) (*types.QueryClaimRecordResponse, error) {

}

func (k Keeper) TotalClaimable(context.Context, *types.QueryTotalClaimableRequest) (*types.QueryTotalClaimableResponse, error) {

}
