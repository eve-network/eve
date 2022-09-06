package keeper

import (
	"github.com/notional-labs/eve/x/tokenfactory/types"
)

var _ types.QueryServer = Keeper{}
