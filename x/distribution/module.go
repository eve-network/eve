package distribution

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	// "github.com/notional-labs/eve/util"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	db "github.com/tendermint/tm-db"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	distr.AppModuleBasic
}

type AppModule struct {
	distr.AppModule
	k  DistributionKeeper
	ak AccountKeeper
	bk bankkeeper.ViewKeeper
	db db.DB
}

func NewAppModule(nested distr.AppModule, k DistributionKeeper, ak AccountKeeper, bk bankkeeper.ViewKeeper, db db.DB) AppModule {
	return AppModule{
		AppModule: nested,
		k:         k,
		ak:        ak,
		bk:        bk,
		db:        db,
	}
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.k, am.ak, am.bk, am.db)
}

// DefaultGenesis returns default genesis state as raw bytes for the distribution module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	state := distrtypes.DefaultGenesisState()
	state.Params.CommunityTax = sdk.ZeroDec()
	return cdc.MustMarshalJSON(state)
}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := cli.NewTxCmd()
	// util.RemoveCobraCommands(cmd, "fund-community-pool") // TODO: ?
	return cmd
}
