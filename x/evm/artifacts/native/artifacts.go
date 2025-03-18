package native

import (
	"bytes"
	"embed"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const CurrentVersion uint16 = 1

//go:embed NativeSheTokensERC20.abi
//go:embed NativeSheTokensERC20.bin
var f embed.FS

var cachedBin []byte
var cachedABI *abi.ABI

func GetABI() []byte {
	bz, err := f.ReadFile("NativeSheTokensERC20.abi")
	if err != nil {
		panic("failed to read NativeSheTokensERC20 contract ABI")
	}
	return bz
}

func GetParsedABI() *abi.ABI {
	if cachedABI != nil {
		return cachedABI
	}
	parsedABI, err := abi.JSON(strings.NewReader(string(GetABI())))
	if err != nil {
		panic(err)
	}
	cachedABI = &parsedABI
	return cachedABI
}

func GetBin() []byte {
	if cachedBin != nil {
		return cachedBin
	}
	code, err := f.ReadFile("NativeSheTokensERC20.bin")
	if err != nil {
		panic("failed to read NativeSheTokensERC20 contract binary")
	}
	bz, err := hex.DecodeString(string(code))
	if err != nil {
		panic("failed to decode NativeSheTokensERC20 contract binary")
	}
	cachedBin = bz
	return bz
}

func IsCodeFromBin(code []byte) bool {
	binLen := len(GetBin())
	if len(code) < binLen {
		return false
	}
	if !bytes.Equal(code[:binLen], GetBin()) {
		return false
	}
	abi, err := NativeMetaData.GetAbi()
	if err != nil {
		fmt.Printf("error getting metadata ABI: %s\n", err)
		return false
	}
	args, err := abi.Constructor.Inputs.Unpack(code[binLen:])
	if err != nil || len(args) != 4 {
		return false
	}
	_, isString1 := args[0].(string)
	_, isString2 := args[1].(string)
	_, isString3 := args[2].(string)
	_, isUint8 := args[3].(uint8)
	return isString1 && isString2 && isString3 && isUint8
}
