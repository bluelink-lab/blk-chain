package wasmbinding

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	aclkeeper "github.com/cosmos/cosmos-sdk/x/accesscontrol/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	epochwasm "github.com/bluelink-lab/blk-chain/x/epoch/client/wasm"
	epochkeeper "github.com/bluelink-lab/blk-chain/x/epoch/keeper"
	evmwasm "github.com/bluelink-lab/blk-chain/x/evm/client/wasm"
	evmkeeper "github.com/bluelink-lab/blk-chain/x/evm/keeper"
	oraclewasm "github.com/bluelink-lab/blk-chain/x/oracle/client/wasm"
	oraclekeeper "github.com/bluelink-lab/blk-chain/x/oracle/keeper"
	tokenfactorywasm "github.com/bluelink-lab/blk-chain/x/tokenfactory/client/wasm"
	tokenfactorykeeper "github.com/bluelink-lab/blk-chain/x/tokenfactory/keeper"
)

func RegisterCustomPlugins(
	oracle *oraclekeeper.Keeper,
	epoch *epochkeeper.Keeper,
	tokenfactory *tokenfactorykeeper.Keeper,
	_ *authkeeper.AccountKeeper,
	router wasmkeeper.MessageRouter,
	channelKeeper wasmtypes.ChannelKeeper,
	capabilityKeeper wasmtypes.CapabilityKeeper,
	bankKeeper wasmtypes.Burner,
	unpacker codectypes.AnyUnpacker,
	portSource wasmtypes.ICS20TransferPortSource,
	aclKeeper aclkeeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
) []wasmkeeper.Option {
	oracleHandler := oraclewasm.NewOracleWasmQueryHandler(oracle)
	epochHandler := epochwasm.NewEpochWasmQueryHandler(epoch)
	tokenfactoryHandler := tokenfactorywasm.NewTokenFactoryWasmQueryHandler(tokenfactory)
	evmHandler := evmwasm.NewEVMQueryHandler(evmKeeper)
	wasmQueryPlugin := NewQueryPlugin(oracleHandler, epochHandler, tokenfactoryHandler, evmHandler, stakingKeeper)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})
	messengerHandlerOpt := wasmkeeper.WithMessageHandler(
		CustomMessageHandler(router, channelKeeper, capabilityKeeper, bankKeeper, evmKeeper, unpacker, portSource, aclKeeper),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerHandlerOpt,
	}
}
