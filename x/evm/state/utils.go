package state

import (
	"encoding/binary"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UsheToSweiMultiplier Fields that were denominated in ushe will be converted to swei (1ushe = 10^12swei)
// for existing Ethereum application (which assumes 18 decimal points) to display properly.
var UsheToSweiMultiplier = big.NewInt(1_000_000_000_000)
var SdkUsheToSweiMultiplier = sdk.NewIntFromBigInt(UsheToSweiMultiplier)

var CoinbaseAddressPrefix = []byte("evm_coinbase")

func GetCoinbaseAddress(txIdx int) sdk.AccAddress {
	txIndexBz := make([]byte, 8)
	binary.BigEndian.PutUint64(txIndexBz, uint64(txIdx))
	return append(CoinbaseAddressPrefix, txIndexBz...)
}

func SplitUsheWeiAmount(amt *big.Int) (sdk.Int, sdk.Int) {
	wei := new(big.Int).Mod(amt, UsheToSweiMultiplier)
	ushe := new(big.Int).Quo(amt, UsheToSweiMultiplier)
	return sdk.NewIntFromBigInt(ushe), sdk.NewIntFromBigInt(wei)
}
