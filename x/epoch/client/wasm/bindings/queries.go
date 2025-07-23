package bindings

import "github.com/bluelink-lab/blk-chain/x/epoch/types"

type SheEpochQuery struct {
	// queries the current Epoch
	Epoch *types.QueryEpochRequest `json:"epoch,omitempty"`
}
