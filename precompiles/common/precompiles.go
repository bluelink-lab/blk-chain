package common

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"math/big"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/bluelink-lab/blk-chain/utils/metrics"
	"github.com/bluelink-lab/blk-chain/x/evm/state"
	"github.com/bluelink-lab/blk-chain/x/evm/types"
)

const UnknownMethodCallGas uint64 = 3000

type Contexter interface {
	Ctx() sdk.Context
}

type StateEVMKeeperGetter interface {
	EVMKeeper() state.EVMKeeper
}

type PrecompileExecutor interface {
	RequiredGas([]byte, *abi.Method) uint64
	Execute(ctx sdk.Context, method *abi.Method, caller common.Address, callingContract common.Address, args []interface{}, value *big.Int, readOnly bool, evm *vm.EVM) ([]byte, error)
}

type Precompile struct {
	abi.ABI
	address  common.Address
	name     string
	executor PrecompileExecutor
}

var _ vm.PrecompiledContract = &Precompile{}

func NewPrecompile(a abi.ABI, executor PrecompileExecutor, address common.Address, name string) *Precompile {
	return &Precompile{ABI: a, executor: executor, address: address, name: name}
}

func (p Precompile) RequiredGas(input []byte) uint64 {
	methodID, err := ExtractMethodID(input)
	if err != nil {
		return UnknownMethodCallGas
	}

	method, err := p.ABI.MethodById(methodID)
	if err != nil {
		// This should never happen since this method is going to fail during Run
		return UnknownMethodCallGas
	}
	return p.executor.RequiredGas(input[4:], method)
}

func (p Precompile) Run(evm *vm.EVM, caller common.Address, callingContract common.Address, input []byte, value *big.Int, readOnly bool, isFromDelegateCall bool) (bz []byte, err error) {
	operation := fmt.Sprintf("%s_unknown", p.name)
	defer func() {
		HandlePrecompileError(err, evm, operation)
		if err != nil {
			bz = []byte(err.Error())
			err = vm.ErrExecutionReverted
		}
	}()
	ctx, method, args, err := p.Prepare(evm, input)
	if err != nil {
		return nil, err
	}

	operation = method.Name
	em := ctx.EventManager()
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	ctx = ctx.WithEVMPrecompileCalledFromDelegateCall(isFromDelegateCall)
	bz, err = p.executor.Execute(ctx, method, caller, callingContract, args, value, readOnly, evm)
	if err != nil {
		return bz, err
	}
	events := ctx.EventManager().Events()
	if len(events) > 0 {
		em.EmitEvents(ctx.EventManager().Events())
	}
	return bz, err
}

func HandlePrecompileError(err error, evm *vm.EVM, operation string) {
	if err != nil {
		evm.StateDB.(*state.DBImpl).SetPrecompileError(err)
		metrics.IncrementErrorMetrics(operation, err)
	}
}

func (p Precompile) Prepare(evm *vm.EVM, input []byte) (sdk.Context, *abi.Method, []interface{}, error) {
	ctxer, ok := evm.StateDB.(Contexter)
	if !ok {
		return sdk.Context{}, nil, nil, errors.New("cannot get context from EVM")
	}
	methodID, err := ExtractMethodID(input)
	if err != nil {
		return sdk.Context{}, nil, nil, err
	}
	method, err := p.ABI.MethodById(methodID)
	if err != nil {
		return sdk.Context{}, nil, nil, err
	}

	argsBz := input[4:]
	args, err := method.Inputs.Unpack(argsBz)
	if err != nil {
		return sdk.Context{}, nil, nil, err
	}

	return ctxer.Ctx(), method, args, nil
}

func (p Precompile) GetABI() abi.ABI {
	return p.ABI
}

func (p Precompile) Address() common.Address {
	return p.address
}

func (p Precompile) GetName() string {
	return p.name
}

func (p Precompile) GetExecutor() PrecompileExecutor {
	return p.executor
}

type DynamicGasPrecompileExecutor interface {
	Execute(ctx sdk.Context, method *abi.Method, caller common.Address, callingContract common.Address, args []interface{}, value *big.Int, readOnly bool, evm *vm.EVM, suppliedGas uint64) (ret []byte, remainingGas uint64, err error)
	EVMKeeper() EVMKeeper
}

type DynamicGasPrecompile struct {
	*Precompile
	executor DynamicGasPrecompileExecutor
}

var _ vm.DynamicGasPrecompiledContract = &DynamicGasPrecompile{}

func NewDynamicGasPrecompile(a abi.ABI, executor DynamicGasPrecompileExecutor, address common.Address, name string) *DynamicGasPrecompile {
	return &DynamicGasPrecompile{Precompile: NewPrecompile(a, nil, address, name), executor: executor}
}

func (d DynamicGasPrecompile) RunAndCalculateGas(evm *vm.EVM, caller common.Address, callingContract common.Address, input []byte, suppliedGas uint64, value *big.Int, hooks *tracing.Hooks, readOnly bool, isFromDelegateCall bool) (ret []byte, remainingGas uint64, err error) {
	operation := fmt.Sprintf("%s_unknown", d.name)
	defer func() {
		HandlePrecompileError(err, evm, operation)
		if err != nil {
			ret = []byte(err.Error())
			err = vm.ErrExecutionReverted
		}
	}()
	ctx, method, args, err := d.Prepare(evm, input)
	if err != nil {
		return nil, 0, err
	}
	gasLimit := d.executor.EVMKeeper().GetCosmosGasLimitFromEVMGas(ctx.WithGasMeter(sdk.NewInfiniteGasMeterWithMultiplier(ctx)), suppliedGas)
	ctx = ctx.WithGasMeter(sdk.NewGasMeterWithMultiplier(ctx, gasLimit))
	operation = method.Name
	em := ctx.EventManager()
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	ctx = ctx.WithEVMPrecompileCalledFromDelegateCall(isFromDelegateCall)
	ret, remainingGas, err = d.executor.Execute(ctx, method, caller, callingContract, args, value, readOnly, evm, suppliedGas)
	if err != nil {
		return ret, remainingGas, err
	}
	events := ctx.EventManager().Events()
	if len(events) > 0 {
		em.EmitEvents(ctx.EventManager().Events())
	}
	return ret, remainingGas, err
}

func (d DynamicGasPrecompile) GetExecutor() DynamicGasPrecompileExecutor {
	return d.executor
}

func ValidateArgsLength(args []interface{}, length int) error {
	if len(args) != length {
		return fmt.Errorf("expected %d arguments but got %d", length, len(args))
	}

	return nil
}

func ValidateNonPayable(value *big.Int) error {
	if value != nil && value.Sign() != 0 {
		return errors.New("sending funds to a non-payable function")
	}

	return nil
}

func HandlePaymentUshe(ctx sdk.Context, precompileAddr sdk.AccAddress, payer sdk.AccAddress, value *big.Int, bankKeeper BankKeeper) (sdk.Coin, error) {
	ublt, wei := state.SplitUsheWeiAmount(value)
	if !wei.IsZero() {
		return sdk.Coin{}, fmt.Errorf("selected precompile function does not allow payment with non-zero wei remainder: received %s", value)
	}
	coin := sdk.NewCoin(sdk.MustGetBaseDenom(), ublt)
	// refund payer because the following precompile logic will debit the payments from payer's account
	// this creates a new event manager to avoid surfacing these as cosmos events
	if err := bankKeeper.SendCoins(ctx.WithEventManager(sdk.NewEventManager()), precompileAddr, payer, sdk.NewCoins(coin)); err != nil {
		return sdk.Coin{}, err
	}
	return coin, nil
}

func HandlePaymentUsheWei(ctx sdk.Context, precompileAddr sdk.AccAddress, payer sdk.AccAddress, value *big.Int, bankKeeper BankKeeper) (sdk.Int, sdk.Int, error) {
	ublt, wei := state.SplitUsheWeiAmount(value)
	// refund payer because the following precompile logic will debit the payments from payer's account
	// this creates a new event manager to avoid surfacing these as cosmos events
	if err := bankKeeper.SendCoinsAndWei(ctx.WithEventManager(sdk.NewEventManager()), precompileAddr, payer, ublt, wei); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	return ublt, wei, nil
}

/*
*
she gas = evm gas * multiplier
she gas price = fee / she gas = fee / (evm gas * multiplier) = evm gas / multiplier
*/
func GetRemainingGas(ctx sdk.Context, evmKeeper EVMKeeper) uint64 {
	return evmKeeper.GetEVMGasLimitFromCtx(ctx)
}

func ExtractMethodID(input []byte) ([]byte, error) {
	// Check if the input has at least the length needed for methodID
	if len(input) < 4 {
		return nil, errors.New("input too short to extract method ID")
	}
	return input[:4], nil
}

func DefaultGasCost(input []byte, isTransaction bool) uint64 {
	if isTransaction {
		return storetypes.KVGasConfig().WriteCostFlat + (storetypes.KVGasConfig().WriteCostPerByte * uint64(len(input)))
	}

	return storetypes.KVGasConfig().ReadCostFlat + (storetypes.KVGasConfig().ReadCostPerByte * uint64(len(input)))
}

func MustGetABI(f embed.FS, filename string) abi.ABI {
	abiBz, err := f.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	newAbi, err := abi.JSON(bytes.NewReader(abiBz))
	if err != nil {
		panic(err)
	}
	return newAbi
}

func GetSheAddressByEvmAddress(ctx sdk.Context, evmAddress common.Address, evmKeeper EVMKeeper) (sdk.AccAddress, error) {
	sheAddr, associated := evmKeeper.GetSheAddress(ctx, evmAddress)
	if !associated {
		return nil, types.NewAssociationMissingErr(evmAddress.Hex())
	}
	return sheAddr, nil
}

func GetSheAddressFromArg(ctx sdk.Context, arg interface{}, evmKeeper EVMKeeper) (sdk.AccAddress, error) {
	addr := arg.(common.Address)
	if addr == (common.Address{}) {
		return nil, errors.New("invalid addr")
	}
	return GetSheAddressByEvmAddress(ctx, addr, evmKeeper)
}
