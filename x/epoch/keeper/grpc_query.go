package keeper

import (
	"github.com/she-protocol/she-chain/x/epoch/types"
)

var _ types.QueryServer = Keeper{}
