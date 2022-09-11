package wasmbinding

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/eve/wasmbinding/bindings"
	tokenfactorykeeper "github.com/notional-labs/eve/x/tokenfactory/keeper"
)

type QueryPlugin struct {
	tokenFactoryKeeper *tokenfactorykeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(tfk *tokenfactorykeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		tokenFactoryKeeper: tfk,
	}
}

// GetDenomAdmin is a query to get denom admin.
func (qp QueryPlugin) GetDenomAdmin(ctx sdk.Context, denom string) (*bindings.DenomAdminResponse, error) {
	metadata, found := qp.tokenFactoryKeeper.GetDenom(ctx, denom)
	if !found {
		return nil, fmt.Errorf("failed to get admin/owner for denom: %s", denom)
	}

	return &bindings.DenomAdminResponse{Admin: metadata.Owner}, nil
}

func (qp QueryPlugin) GetDenomData(ctx sdk.Context, denom string) (*bindings.FullDenomResponse, error) {
	data, found := qp.tokenFactoryKeeper.GetDenom(ctx, denom)
	if !found {
		return nil, fmt.Errorf("failed to get admin/owner for denom: %s", denom)
	}

	return &bindings.FullDenomResponse{
		Name:               data.Name,
		Denom:              data.Denom,
		Precision:          data.Precision,
		MaxSupply:          data.MaxSupply,
		Supply:             data.Supply,
		CanChangeMaxSupply: data.CanChangeMaxSupply,
		Owner:              data.Owner,
		// metadata
		Description: data.Description,
		TokenImage:  data.TokenImage,
		Website:     data.Website,
	}, nil
}
