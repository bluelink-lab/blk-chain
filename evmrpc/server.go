package evmrpc

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	evmCfg "github.com/bluelink-lab/blk-chain/x/evm/config"
	"github.com/bluelink-lab/blk-chain/x/evm/keeper"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type ConnectionType string

var ConnectionTypeWS ConnectionType = "websocket"
var ConnectionTypeHTTP ConnectionType = "http"

const LocalAddress = "0.0.0.0"

type EVMServer interface {
	Start() error
	Stop()
}

func NewEVMHTTPServer(
	logger log.Logger,
	config Config,
	tmClient rpcclient.Client,
	k *keeper.Keeper,
	app *baseapp.BaseApp,
	antehandler sdk.AnteHandler,
	ctxProvider func(int64) sdk.Context,
	txConfig client.TxConfig,
	homeDir string,
	isPanicOrSyntheticTxFunc func(ctx context.Context, hash common.Hash) (bool, error), // used in *ExcludeTraceFail endpoints
) (EVMServer, error) {
	httpServer := NewHTTPServer(logger, rpc.HTTPTimeouts{
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
	})
	if err := httpServer.SetListenAddr(LocalAddress, config.HTTPPort); err != nil {
		return nil, err
	}
	simulateConfig := &SimulateConfig{GasCap: config.SimulationGasLimit, EVMTimeout: config.SimulationEVMTimeout}
	sendAPI := NewSendAPI(tmClient, txConfig, &SendConfig{slow: config.Slow}, k, ctxProvider, homeDir, simulateConfig, app, antehandler, ConnectionTypeHTTP)
	ctx := ctxProvider(LatestCtxHeight)

	txAPI := NewTransactionAPI(tmClient, k, ctxProvider, txConfig, homeDir, ConnectionTypeHTTP)
	debugAPI := NewDebugAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), simulateConfig, app, antehandler, ConnectionTypeHTTP)
	if isPanicOrSyntheticTxFunc == nil {
		isPanicOrSyntheticTxFunc = func(ctx context.Context, hash common.Hash) (bool, error) {
			return debugAPI.isPanicOrSyntheticTx(ctx, hash)
		}
	}
	sheTxAPI := NewSheTransactionAPI(tmClient, k, ctxProvider, txConfig, homeDir, ConnectionTypeHTTP, isPanicOrSyntheticTxFunc)
	sheDebugAPI := NewSheDebugAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), simulateConfig, app, antehandler, ConnectionTypeHTTP)

	apis := []rpc.API{
		{
			Namespace: "echo",
			Service:   NewEchoAPI(),
		},
		{
			Namespace: "eth",
			Service:   NewBlockAPI(tmClient, k, ctxProvider, txConfig, ConnectionTypeHTTP),
		},
		{
			Namespace: "blt",
			Service:   NewSheBlockAPI(tmClient, k, ctxProvider, txConfig, ConnectionTypeHTTP, isPanicOrSyntheticTxFunc),
		},
		{
			Namespace: "blt2",
			Service:   NewShe2BlockAPI(tmClient, k, ctxProvider, txConfig, ConnectionTypeHTTP, isPanicOrSyntheticTxFunc),
		},
		{
			Namespace: "eth",
			Service:   txAPI,
		},
		{
			Namespace: "blt",
			Service:   sheTxAPI,
		},
		{
			Namespace: "eth",
			Service:   NewStateAPI(tmClient, k, ctxProvider, ConnectionTypeHTTP),
		},
		{
			Namespace: "eth",
			Service:   NewInfoAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), homeDir, config.MaxBlocksForLog, ConnectionTypeHTTP),
		},
		{
			Namespace: "eth",
			Service:   sendAPI,
		},
		{
			Namespace: "eth",
			Service:   NewSimulationAPI(ctxProvider, k, txConfig.TxDecoder(), tmClient, simulateConfig, app, antehandler, ConnectionTypeHTTP),
		},
		{
			Namespace: "net",
			Service:   NewNetAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), ConnectionTypeHTTP),
		},
		{
			Namespace: "eth",
			Service:   NewFilterAPI(tmClient, k, ctxProvider, txConfig, &FilterConfig{timeout: config.FilterTimeout, maxLog: config.MaxLogNoBlock, maxBlock: config.MaxBlocksForLog}, ConnectionTypeHTTP, "eth"),
		},
		{
			Namespace: "blt",
			Service:   NewFilterAPI(tmClient, k, ctxProvider, txConfig, &FilterConfig{timeout: config.FilterTimeout, maxLog: config.MaxLogNoBlock, maxBlock: config.MaxBlocksForLog}, ConnectionTypeHTTP, "blt"),
		},
		{
			Namespace: "blt",
			Service:   NewAssociationAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), sendAPI, ConnectionTypeHTTP),
		},
		{
			Namespace: "txpool",
			Service:   NewTxPoolAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), &TxPoolConfig{maxNumTxs: int(config.MaxTxPoolTxs)}, ConnectionTypeHTTP),
		},
		{
			Namespace: "web3",
			Service:   &Web3API{},
		},
		{
			Namespace: "debug",
			Service:   debugAPI,
		},
		{
			Namespace: "blt",
			Service:   sheDebugAPI,
		},
	}
	// Test API can only exist on non-live chain IDs.  These APIs instrument certain overrides.
	if config.EnableTestAPI && !evmCfg.IsLiveChainID(ctx) {
		logger.Info("Enabling Test EVM APIs")
		apis = append(apis, rpc.API{
			Namespace: "test",
			Service:   NewTestAPI(),
		})
	} else {
		logger.Info("Disabling Test EVM APIs", "liveChainID", evmCfg.IsLiveChainID(ctx), "enableTestAPI", config.EnableTestAPI)
	}

	if err := httpServer.EnableRPC(apis, HTTPConfig{
		CorsAllowedOrigins: strings.Split(config.CORSOrigins, ","),
		Vhosts:             []string{"*"},
	}); err != nil {
		return nil, err
	}
	return httpServer, nil
}

func NewEVMWebSocketServer(
	logger log.Logger,
	config Config,
	tmClient rpcclient.Client,
	k *keeper.Keeper,
	app *baseapp.BaseApp,
	antehandler sdk.AnteHandler,
	ctxProvider func(int64) sdk.Context,
	txConfig client.TxConfig,
	homeDir string,
) (EVMServer, error) {
	httpServer := NewHTTPServer(logger, rpc.HTTPTimeouts{
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
	})
	if err := httpServer.SetListenAddr(LocalAddress, config.WSPort); err != nil {
		return nil, err
	}
	simulateConfig := &SimulateConfig{GasCap: config.SimulationGasLimit, EVMTimeout: config.SimulationEVMTimeout}
	apis := []rpc.API{
		{
			Namespace: "echo",
			Service:   NewEchoAPI(),
		},
		{
			Namespace: "eth",
			Service:   NewBlockAPI(tmClient, k, ctxProvider, txConfig, ConnectionTypeWS),
		},
		{
			Namespace: "eth",
			Service:   NewTransactionAPI(tmClient, k, ctxProvider, txConfig, homeDir, ConnectionTypeWS),
		},
		{
			Namespace: "eth",
			Service:   NewStateAPI(tmClient, k, ctxProvider, ConnectionTypeWS),
		},
		{
			Namespace: "eth",
			Service:   NewInfoAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), homeDir, config.MaxBlocksForLog, ConnectionTypeWS),
		},
		{
			Namespace: "eth",
			Service:   NewSendAPI(tmClient, txConfig, &SendConfig{slow: config.Slow}, k, ctxProvider, homeDir, simulateConfig, app, antehandler, ConnectionTypeWS),
		},
		{
			Namespace: "eth",
			Service:   NewSimulationAPI(ctxProvider, k, txConfig.TxDecoder(), tmClient, simulateConfig, app, antehandler, ConnectionTypeWS),
		},
		{
			Namespace: "net",
			Service:   NewNetAPI(tmClient, k, ctxProvider, txConfig.TxDecoder(), ConnectionTypeWS),
		},
		{
			Namespace: "eth",
			Service:   NewSubscriptionAPI(tmClient, k, ctxProvider, &LogFetcher{tmClient: tmClient, k: k, ctxProvider: ctxProvider, txConfig: txConfig}, &SubscriptionConfig{subscriptionCapacity: 100, newHeadLimit: config.MaxSubscriptionsNewHead}, &FilterConfig{timeout: config.FilterTimeout, maxLog: config.MaxLogNoBlock, maxBlock: config.MaxBlocksForLog}, ConnectionTypeWS),
		},
		{
			Namespace: "web3",
			Service:   &Web3API{},
		},
	}
	if err := httpServer.EnableWS(apis, WsConfig{Origins: strings.Split(config.WSOrigins, ",")}); err != nil {
		return nil, err
	}
	return httpServer, nil
}
