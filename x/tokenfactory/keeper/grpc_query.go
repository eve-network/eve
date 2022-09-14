package keeper

import (
	"github.com/eve-network/eve/x/tokenfactory/types"
)

var _ types.QueryServer = Keeper{}
