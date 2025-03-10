package addr

import (
	"bytes"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"

	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/she-protocol/she-chain/utils"
	"github.com/she-protocol/she-chain/utils/helpers"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	pcommon "github.com/she-protocol/she-chain/precompiles/common"
	"github.com/she-protocol/she-chain/utils/metrics"
	"github.com/she-protocol/she-chain/x/evm/types"
)

const (
	GetSheAddressMethod = "getSheAddr"
	GetEvmAddressMethod = "getEvmAddr"
	Associate           = "associate"
	AssociatePubKey     = "associatePubKey"
)

const (
	AddrAddress = "0x0000000000000000000000000000000000001004"
)

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

type PrecompileExecutor struct {
	evmKeeper     pcommon.EVMKeeper
	bankKeeper    pcommon.BankKeeper
	accountKeeper pcommon.AccountKeeper

	GetSheAddressID   []byte
	GetEvmAddressID   []byte
	AssociateID       []byte
	AssociatePubKeyID []byte
}

func NewPrecompile(evmKeeper pcommon.EVMKeeper, bankKeeper pcommon.BankKeeper, accountKeeper pcommon.AccountKeeper) (*pcommon.DynamicGasPrecompile, error) {

	newAbi := pcommon.MustGetABI(f, "abi.json")

	p := &PrecompileExecutor{
		evmKeeper:     evmKeeper,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}

	for name, m := range newAbi.Methods {
		switch name {
		case GetSheAddressMethod:
			p.GetSheAddressID = m.ID
		case GetEvmAddressMethod:
			p.GetEvmAddressID = m.ID
		case Associate:
			p.AssociateID = m.ID
		case AssociatePubKey:
			p.AssociatePubKeyID = m.ID
		}
	}

	return pcommon.NewDynamicGasPrecompile(newAbi, p, common.HexToAddress(AddrAddress), "addr"), nil
}

// RequiredGas returns the required bare minimum gas to execute the precompile.
func (p PrecompileExecutor) RequiredGas(input []byte, method *abi.Method) uint64 {
	if bytes.Equal(method.ID, p.AssociateID) || bytes.Equal(method.ID, p.AssociatePubKeyID) {
		return 50000
	}
	return pcommon.DefaultGasCost(input, p.IsTransaction(method.Name))
}

func (p PrecompileExecutor) Execute(ctx sdk.Context, method *abi.Method, _ common.Address, _ common.Address, args []interface{}, value *big.Int, readOnly bool, _ *vm.EVM, suppliedGas uint64) (ret []byte, remainingGas uint64, err error) {
	// Needed to catch gas meter panics
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("execution reverted: %v", r)
		}
	}()
	switch method.Name {
	case GetSheAddressMethod:
		return p.getSheAddr(ctx, method, args, value)
	case GetEvmAddressMethod:
		return p.getEvmAddr(ctx, method, args, value)
	case Associate:
		if readOnly {
			return nil, 0, errors.New("cannot call associate precompile from staticcall")
		}
		return p.associate(ctx, method, args, value)
	case AssociatePubKey:
		if readOnly {
			return nil, 0, errors.New("cannot call associate pub key precompile from staticcall")
		}
		return p.associatePublicKey(ctx, method, args, value)
	}
	return
}

func (p PrecompileExecutor) getSheAddr(ctx sdk.Context, method *abi.Method, args []interface{}, value *big.Int) (ret []byte, remainingGas uint64, err error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, 0, err
	}

	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, 0, err
	}

	sheAddr, found := p.evmKeeper.GetSheAddress(ctx, args[0].(common.Address))
	if !found {
		metrics.IncrementAssociationError("getSheAddr", types.NewAssociationMissingErr(args[0].(common.Address).Hex()))
		return nil, 0, fmt.Errorf("EVM address %s is not associated", args[0].(common.Address).Hex())
	}
	ret, err = method.Outputs.Pack(sheAddr.String())
	return ret, pcommon.GetRemainingGas(ctx, p.evmKeeper), err
}

func (p PrecompileExecutor) getEvmAddr(ctx sdk.Context, method *abi.Method, args []interface{}, value *big.Int) (ret []byte, remainingGas uint64, err error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, 0, err
	}

	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, 0, err
	}

	sheAddr, err := sdk.AccAddressFromBech32(args[0].(string))
	if err != nil {
		return nil, 0, err
	}

	evmAddr, found := p.evmKeeper.GetEVMAddress(ctx, sheAddr)
	if !found {
		metrics.IncrementAssociationError("getEvmAddr", types.NewAssociationMissingErr(args[0].(string)))
		return nil, 0, fmt.Errorf("she address %s is not associated", args[0].(string))
	}
	ret, err = method.Outputs.Pack(evmAddr)
	return ret, pcommon.GetRemainingGas(ctx, p.evmKeeper), err
}

func (p PrecompileExecutor) associate(ctx sdk.Context, method *abi.Method, args []interface{}, value *big.Int) (ret []byte, remainingGas uint64, err error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, 0, err
	}

	if err := pcommon.ValidateArgsLength(args, 4); err != nil {
		return nil, 0, err
	}

	// v, r and s are components of a signature over the customMessage sent.
	// We use the signature to construct the user's pubkey to obtain their addresses.
	v := args[0].(string)
	r := args[1].(string)
	s := args[2].(string)
	customMessage := args[3].(string)

	rBytes, err := decodeHexString(r)
	if err != nil {
		return nil, 0, err
	}
	sBytes, err := decodeHexString(s)
	if err != nil {
		return nil, 0, err
	}
	vBytes, err := decodeHexString(v)
	if err != nil {
		return nil, 0, err
	}

	vBig := new(big.Int).SetBytes(vBytes)
	rBig := new(big.Int).SetBytes(rBytes)
	sBig := new(big.Int).SetBytes(sBytes)

	// Derive addresses
	vBig = new(big.Int).Add(vBig, utils.Big27)

	customMessageHash := crypto.Keccak256Hash([]byte(customMessage))
	evmAddr, sheAddr, pubkey, err := helpers.GetAddresses(vBig, rBig, sBig, customMessageHash)
	if err != nil {
		return nil, 0, err
	}

	return p.associateAddresses(ctx, method, evmAddr, sheAddr, pubkey)
}

func (p PrecompileExecutor) associatePublicKey(ctx sdk.Context, method *abi.Method, args []interface{}, value *big.Int) (ret []byte, remainingGas uint64, err error) {
	if err := pcommon.ValidateNonPayable(value); err != nil {
		return nil, 0, err
	}

	if err := pcommon.ValidateArgsLength(args, 1); err != nil {
		return nil, 0, err
	}

	// Takes a single argument, a compressed pubkey in hex format, excluding the '0x'
	pubKeyHex := args[0].(string)

	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, 0, err
	}

	// Parse the compressed public key
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return nil, 0, err
	}

	// Convert to uncompressed public key
	uncompressedPubKey := pubKey.SerializeUncompressed()

	evmAddr, sheAddr, pubkey, err := helpers.GetAddressesFromPubkeyBytes(uncompressedPubKey)
	if err != nil {
		return nil, 0, err
	}

	return p.associateAddresses(ctx, method, evmAddr, sheAddr, pubkey)
}

func (p PrecompileExecutor) associateAddresses(ctx sdk.Context, method *abi.Method, evmAddr common.Address, sheAddr sdk.AccAddress, pubkey cryptotypes.PubKey) (ret []byte, remainingGas uint64, err error) {
	// Check that address is not already associated
	_, found := p.evmKeeper.GetEVMAddress(ctx, sheAddr)
	if found {
		return nil, 0, fmt.Errorf("address %s is already associated with evm address %s", sheAddr, evmAddr)
	}

	// Associate Addresses:
	associationHelper := helpers.NewAssociationHelper(p.evmKeeper, p.bankKeeper, p.accountKeeper)
	err = associationHelper.AssociateAddresses(ctx, sheAddr, evmAddr, pubkey)
	if err != nil {
		return nil, 0, err
	}

	ret, err = method.Outputs.Pack(sheAddr.String(), evmAddr)
	return ret, pcommon.GetRemainingGas(ctx, p.evmKeeper), err
}

func (PrecompileExecutor) IsTransaction(method string) bool {
	switch method {
	case Associate:
		return true
	default:
		return false
	}
}

func (p PrecompileExecutor) EVMKeeper() pcommon.EVMKeeper {
	return p.evmKeeper
}

func decodeHexString(hexString string) ([]byte, error) {
	trimmed := strings.TrimPrefix(hexString, "0x")
	if len(trimmed)%2 != 0 {
		trimmed = "0" + trimmed
	}
	return hex.DecodeString(trimmed)
}
