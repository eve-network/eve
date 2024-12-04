package app

import (
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v9/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

func (app *EveApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *EveApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *EveApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *EveApp) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

func (app *EveApp) GetStakingKeeper() *stakingkeeper.Keeper {
	return &app.StakingKeeper
}

func (app *EveApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *EveApp) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}
