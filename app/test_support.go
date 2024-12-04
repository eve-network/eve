package app

import (
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

func (app *LimeApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *LimeApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *LimeApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *LimeApp) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

func (app *LimeApp) GetStakingKeeper() *stakingkeeper.Keeper {
	return &app.StakingKeeper
}

func (app *LimeApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *LimeApp) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}
