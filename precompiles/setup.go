package precompiles

import (
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/bluelink-lab/blk-chain/precompiles/addr"
	"github.com/bluelink-lab/blk-chain/precompiles/bank"
	"github.com/bluelink-lab/blk-chain/precompiles/common"
	"github.com/bluelink-lab/blk-chain/precompiles/distribution"
	"github.com/bluelink-lab/blk-chain/precompiles/gov"
	"github.com/bluelink-lab/blk-chain/precompiles/ibc"
	"github.com/bluelink-lab/blk-chain/precompiles/json"
	"github.com/bluelink-lab/blk-chain/precompiles/oracle"
	"github.com/bluelink-lab/blk-chain/precompiles/pointer"
	"github.com/bluelink-lab/blk-chain/precompiles/pointerview"
	"github.com/bluelink-lab/blk-chain/precompiles/staking"
	"github.com/bluelink-lab/blk-chain/precompiles/wasmd"
)

var SetupMtx = &sync.Mutex{}
var Initialized = false

type PrecompileInfo struct {
	ABI     abi.ABI
	Address ecommon.Address
}

// PrecompileNamesToInfo is Populated by InitializePrecompiles
var PrecompileNamesToInfo = map[string]PrecompileInfo{}

type IPrecompile interface {
	vm.PrecompiledContract
	GetABI() abi.ABI
	GetName() string
	Address() ecommon.Address
}

func GetCustomPrecompiles(
	evmKeeper common.EVMKeeper,
	bankKeeper common.BankKeeper,
	bankSender common.BankMsgServer,
	wasmdKeeper common.WasmdKeeper,
	wasmdViewKeeper common.WasmdViewKeeper,
	stakingKeeper common.StakingKeeper,
	stakingQuerier common.StakingQuerier,
	govKeeper common.GovKeeper,
	distrKeeper common.DistributionKeeper,
	oracleKeeper common.OracleKeeper,
	transferKeeper common.TransferKeeper,
	clientKeeper common.ClientKeeper,
	connectionKeeper common.ConnectionKeeper,
	channelKeeper common.ChannelKeeper,
	accountKeeper common.AccountKeeper,
) map[ecommon.Address]vm.PrecompiledContract {
	bankp, err := bank.NewPrecompile(bankKeeper, bankSender, evmKeeper, accountKeeper)
	if err != nil {
		panic(err)
	}
	wasmdp, err := wasmd.NewPrecompile(evmKeeper, wasmdKeeper, wasmdViewKeeper, bankKeeper)
	if err != nil {
		panic(err)
	}
	jsonp, err := json.NewPrecompile()
	if err != nil {
		panic(err)
	}
	addrp, err := addr.NewPrecompile(evmKeeper, bankKeeper, accountKeeper)
	if err != nil {
		panic(err)
	}
	stakingp, err := staking.NewPrecompile(stakingKeeper, stakingQuerier, evmKeeper, bankKeeper)
	if err != nil {
		panic(err)
	}
	govp, err := gov.NewPrecompile(govKeeper, evmKeeper, bankKeeper)
	if err != nil {
		panic(err)
	}
	distrp, err := distribution.NewPrecompile(distrKeeper, evmKeeper)
	if err != nil {
		panic(err)
	}
	oraclep, err := oracle.NewPrecompile(oracleKeeper, evmKeeper)
	if err != nil {
		panic(err)
	}
	ibcp, err := ibc.NewPrecompile(transferKeeper, evmKeeper, clientKeeper, connectionKeeper, channelKeeper)
	if err != nil {
		panic(err)
	}
	pointerp, err := pointer.NewPrecompile(evmKeeper, bankKeeper, wasmdViewKeeper)
	if err != nil {
		panic(err)
	}
	pointerviewp, err := pointerview.NewPrecompile(evmKeeper)
	if err != nil {
		panic(err)
	}
	return map[ecommon.Address]vm.PrecompiledContract{
		bankp.Address():        bankp,
		wasmdp.Address():       wasmdp,
		jsonp.Address():        jsonp,
		addrp.Address():        addrp,
		stakingp.Address():     stakingp,
		govp.Address():         govp,
		distrp.Address():       distrp,
		oraclep.Address():      oraclep,
		ibcp.Address():         ibcp,
		pointerp.Address():     pointerp,
		pointerviewp.Address(): pointerviewp,
	}
}

func InitializePrecompiles(
	dryRun bool,
	evmKeeper common.EVMKeeper,
	bankKeeper common.BankKeeper,
	bankSender common.BankMsgServer,
	wasmdKeeper common.WasmdKeeper,
	wasmdViewKeeper common.WasmdViewKeeper,
	stakingKeeper common.StakingKeeper,
	stakingQuerier common.StakingQuerier,
	govKeeper common.GovKeeper,
	distrKeeper common.DistributionKeeper,
	oracleKeeper common.OracleKeeper,
	transferKeeper common.TransferKeeper,
	clientKeeper common.ClientKeeper,
	connectionKeeper common.ConnectionKeeper,
	channelKeeper common.ChannelKeeper,
	accountKeeper common.AccountKeeper,
) error {
	SetupMtx.Lock()
	defer SetupMtx.Unlock()
	if Initialized {
		panic("precompiles already initialized")
	}
	bankp, err := bank.NewPrecompile(bankKeeper, bankSender, evmKeeper, accountKeeper)
	if err != nil {
		return err
	}
	wasmdp, err := wasmd.NewPrecompile(evmKeeper, wasmdKeeper, wasmdViewKeeper, bankKeeper)
	if err != nil {
		return err
	}
	jsonp, err := json.NewPrecompile()
	if err != nil {
		return err
	}
	addrp, err := addr.NewPrecompile(evmKeeper, bankKeeper, accountKeeper)
	if err != nil {
		return err
	}
	stakingp, err := staking.NewPrecompile(stakingKeeper, stakingQuerier, evmKeeper, bankKeeper)
	if err != nil {
		return err
	}
	govp, err := gov.NewPrecompile(govKeeper, evmKeeper, bankKeeper)
	if err != nil {
		return err
	}
	distrp, err := distribution.NewPrecompile(distrKeeper, evmKeeper)
	if err != nil {
		return err
	}
	oraclep, err := oracle.NewPrecompile(oracleKeeper, evmKeeper)
	if err != nil {
		return err
	}
	ibcp, err := ibc.NewPrecompile(transferKeeper, evmKeeper, clientKeeper, connectionKeeper, channelKeeper)
	if err != nil {
		return err
	}
	pointerp, err := pointer.NewPrecompile(evmKeeper, bankKeeper, wasmdViewKeeper)
	if err != nil {
		return err
	}
	pointerviewp, err := pointerview.NewPrecompile(evmKeeper)
	if err != nil {
		return err
	}
	PrecompileNamesToInfo[bankp.GetName()] = PrecompileInfo{ABI: bankp.GetABI(), Address: bankp.Address()}
	PrecompileNamesToInfo[wasmdp.GetName()] = PrecompileInfo{ABI: wasmdp.GetABI(), Address: wasmdp.Address()}
	PrecompileNamesToInfo[jsonp.GetName()] = PrecompileInfo{ABI: jsonp.GetABI(), Address: jsonp.Address()}
	PrecompileNamesToInfo[addrp.GetName()] = PrecompileInfo{ABI: addrp.GetABI(), Address: addrp.Address()}
	PrecompileNamesToInfo[stakingp.GetName()] = PrecompileInfo{ABI: stakingp.GetABI(), Address: stakingp.Address()}
	PrecompileNamesToInfo[govp.GetName()] = PrecompileInfo{ABI: govp.GetABI(), Address: govp.Address()}
	PrecompileNamesToInfo[distrp.GetName()] = PrecompileInfo{ABI: distrp.GetABI(), Address: distrp.Address()}
	PrecompileNamesToInfo[oraclep.GetName()] = PrecompileInfo{ABI: oraclep.GetABI(), Address: oraclep.Address()}
	PrecompileNamesToInfo[ibcp.GetName()] = PrecompileInfo{ABI: ibcp.GetABI(), Address: ibcp.Address()}
	PrecompileNamesToInfo[pointerp.GetName()] = PrecompileInfo{ABI: pointerp.GetABI(), Address: pointerp.Address()}
	PrecompileNamesToInfo[pointerviewp.GetName()] = PrecompileInfo{ABI: pointerviewp.GetABI(), Address: pointerviewp.Address()}
	if !dryRun {
		addPrecompileToVM(bankp)
		addPrecompileToVM(wasmdp)
		addPrecompileToVM(jsonp)
		addPrecompileToVM(addrp)
		addPrecompileToVM(stakingp)
		addPrecompileToVM(govp)
		addPrecompileToVM(distrp)
		addPrecompileToVM(oraclep)
		addPrecompileToVM(ibcp)
		addPrecompileToVM(pointerp)
		addPrecompileToVM(pointerviewp)
		Initialized = true
	}
	return nil
}

func GetPrecompileInfo(name string) PrecompileInfo {
	if !Initialized {
		// Precompile Info does not require any keeper state
		_ = InitializePrecompiles(true, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}
	i, ok := PrecompileNamesToInfo[name]
	if !ok {
		panic(name + "doesn't exist as a precompile")
	}
	return i
}

// This function modifies global variable in `vm` module. It should only be called once
// per precompile during initialization
func addPrecompileToVM(p IPrecompile) {
	vm.PrecompiledContractsHomestead[p.Address()] = p
	vm.PrecompiledContractsByzantium[p.Address()] = p
	vm.PrecompiledContractsIstanbul[p.Address()] = p
	vm.PrecompiledContractsBerlin[p.Address()] = p
	vm.PrecompiledContractsCancun[p.Address()] = p
	vm.PrecompiledContractsBLS[p.Address()] = p
	vm.PrecompiledAddressesHomestead = append(vm.PrecompiledAddressesHomestead, p.Address())
	vm.PrecompiledAddressesByzantium = append(vm.PrecompiledAddressesByzantium, p.Address())
	vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, p.Address())
	vm.PrecompiledAddressesBerlin = append(vm.PrecompiledAddressesBerlin, p.Address())
	vm.PrecompiledAddressesCancun = append(vm.PrecompiledAddressesCancun, p.Address())
}
