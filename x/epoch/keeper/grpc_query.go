package keeper

import (
	"github.com/bluelink-lab/blk-chain/x/epoch/types"
)

var _ types.QueryServer = Keeper{}
