package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/epoch module sentinel errors
var (
	ErrParsingSheEpochQuery = sdkerrors.Register(ModuleName, 2, "Error parsing SheEpochQuery")
	ErrGettingEpoch         = sdkerrors.Register(ModuleName, 3, "Error while getting epoch")
	ErrEncodingEpoch        = sdkerrors.Register(ModuleName, 4, "Error encoding epoch as JSON")
	ErrUnknownSheEpochQuery = sdkerrors.Register(ModuleName, 6, "Error unknown she epoch query")
)
