package bindings

import "github.com/she-protocol/she-chain/x/epoch/types"

type SheEpochQuery struct {
	// queries the current Epoch
	Epoch *types.QueryEpochRequest `json:"epoch,omitempty"`
}
