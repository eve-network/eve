package v1

import (
	"context"

	"github.com/eve-network/eve/app/upgrades"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// TODO: Add any additional upgrades here
func CreateUpgradeHandler(mm upgrades.ModuleManager,
	configurator module.Configurator,
	keepers *upgrades.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func ConfigureFeeMarketModule(ctx sdk.Context, keepers *upgrades.AppKeepers) error {
	params, err := keepers.FeeMarketKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	params.Enabled = true
	params.FeeDenom = "ueve"
	params.DistributeFees = false // burn fees
	params.MinBaseGasPrice = sdkmath.LegacyMustNewDecFromStr("0.005")
	params.MaxBlockUtilization = feemarkettypes.DefaultMaxBlockUtilization
	if err := keepers.FeeMarketKeeper.SetParams(ctx, params); err != nil {
		return err
	}

	state, err := keepers.FeeMarketKeeper.GetState(ctx)
	if err != nil {
		return err
	}

	state.BaseGasPrice = sdkmath.LegacyMustNewDecFromStr("0.005")

	return keepers.FeeMarketKeeper.SetState(ctx, state)
}
