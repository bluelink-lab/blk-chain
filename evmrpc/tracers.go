package evmrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/tracers"
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"     // run init()s to register JS tracers
	_ "github.com/ethereum/go-ethereum/eth/tracers/native" // run init()s to register native tracers
	"github.com/ethereum/go-ethereum/lib/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/bluelink-lab/blk-chain/x/evm/keeper"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

const (
	IsPanicCacheSize = 5000
	IsPanicCacheTTL  = 1 * time.Minute
)

type DebugAPI struct {
	tracersAPI     *tracers.API
	tmClient       rpcclient.Client
	keeper         *keeper.Keeper
	ctxProvider    func(int64) sdk.Context
	txDecoder      sdk.TxDecoder
	connectionType ConnectionType
	isPanicCache   *expirable.LRU[common.Hash, bool] // hash to isPanic
}

type SheDebugAPI struct {
	*DebugAPI
}

func NewDebugAPI(tmClient rpcclient.Client, k *keeper.Keeper, ctxProvider func(int64) sdk.Context, txDecoder sdk.TxDecoder, config *SimulateConfig, app *baseapp.BaseApp,
	antehandler sdk.AnteHandler, connectionType ConnectionType) *DebugAPI {
	backend := NewBackend(ctxProvider, k, txDecoder, tmClient, config, app, antehandler)
	tracersAPI := tracers.NewAPI(backend)
	evictCallback := func(key common.Hash, value bool) {}
	isPanicCache := expirable.NewLRU[common.Hash, bool](IsPanicCacheSize, evictCallback, IsPanicCacheTTL)
	return &DebugAPI{
		tracersAPI:     tracersAPI,
		tmClient:       tmClient,
		keeper:         k,
		ctxProvider:    ctxProvider,
		txDecoder:      txDecoder,
		connectionType: connectionType,
		isPanicCache:   isPanicCache,
	}
}

func NewSheDebugAPI(
	tmClient rpcclient.Client,
	k *keeper.Keeper,
	ctxProvider func(int64) sdk.Context,
	txDecoder sdk.TxDecoder,
	config *SimulateConfig,
	app *baseapp.BaseApp,
	antehandler sdk.AnteHandler,
	connectionType ConnectionType,
) *SheDebugAPI {
	backend := NewBackend(ctxProvider, k, txDecoder, tmClient, config, app, antehandler)
	tracersAPI := tracers.NewAPI(backend)
	return &SheDebugAPI{
		DebugAPI: &DebugAPI{tracersAPI: tracersAPI, tmClient: tmClient, keeper: k, ctxProvider: ctxProvider, txDecoder: txDecoder, connectionType: connectionType},
	}
}

func (api *DebugAPI) TraceTransaction(ctx context.Context, hash common.Hash, config *tracers.TraceConfig) (result interface{}, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("debug_traceTransaction", api.connectionType, startTime, returnErr == nil)
	result, returnErr = api.tracersAPI.TraceTransaction(ctx, hash, config)
	return
}

func (api *SheDebugAPI) TraceBlockByNumberExcludeTraceFail(ctx context.Context, number rpc.BlockNumber, config *tracers.TraceConfig) (result interface{}, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("blt_traceBlockByNumberExcludeTraceFail", api.connectionType, startTime, returnErr == nil)
	result, returnErr = api.tracersAPI.TraceBlockByNumber(ctx, number, config)
	traces, ok := result.([]*tracers.TxTraceResult)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", result)
	}
	finalTraces := make([]*tracers.TxTraceResult, 0, len(traces))
	for _, trace := range traces {
		if len(trace.Error) > 0 {
			continue
		}
		finalTraces = append(finalTraces, trace)
	}
	return finalTraces, nil
}

func (api *SheDebugAPI) TraceBlockByHashExcludeTraceFail(ctx context.Context, hash common.Hash, config *tracers.TraceConfig) (result interface{}, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("blt_traceBlockByHashExcludeTraceFail", api.connectionType, startTime, returnErr == nil)
	result, returnErr = api.tracersAPI.TraceBlockByHash(ctx, hash, config)
	traces, ok := result.([]*tracers.TxTraceResult)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", result)
	}
	finalTraces := make([]*tracers.TxTraceResult, 0, len(traces))
	for _, trace := range traces {
		if len(trace.Error) > 0 {
			continue
		}
		finalTraces = append(finalTraces, trace)
	}
	return finalTraces, nil
}

// isPanicOrSyntheticTx returns true if the tx is a panic tx or if it is a synthetic tx. Used in the *ExcludeTraceFail endpoints.
func (api *DebugAPI) isPanicOrSyntheticTx(ctx context.Context, hash common.Hash) (isPanic bool, err error) {
	sdkctx := api.ctxProvider(LatestCtxHeight)
	receipt, err := api.keeper.GetReceipt(sdkctx, hash)
	if err != nil {
		return false, err
	}
	height := receipt.BlockNumber

	isPanic, ok := api.isPanicCache.Get(hash)
	if ok {
		return isPanic, nil
	}

	callTracer := "callTracer"
	tracersResult, err := api.tracersAPI.TraceBlockByNumber(ctx, rpc.BlockNumber(height), &tracers.TraceConfig{
		Tracer: &callTracer,
	})
	if err != nil {
		return false, err
	}

	found := false
	result := false
	for _, trace := range tracersResult {
		if trace.TxHash == hash {
			found = true
			result = len(trace.Error) > 0
		}
		// for each tx, add to cache to avoid re-tracing
		if len(trace.Error) > 0 {
			api.isPanicCache.Add(trace.TxHash, true)
		} else {
			api.isPanicCache.Add(trace.TxHash, false)
		}
	}

	if !found { // likely a synthetic tx
		return true, nil
	}

	return result, nil
}

func (api *DebugAPI) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *tracers.TraceConfig) (result interface{}, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("debug_traceBlockByNumber", api.connectionType, startTime, returnErr == nil)
	result, returnErr = api.tracersAPI.TraceBlockByNumber(ctx, number, config)
	return
}

func (api *DebugAPI) TraceBlockByHash(ctx context.Context, hash common.Hash, config *tracers.TraceConfig) (result interface{}, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("debug_traceBlockByHash", api.connectionType, startTime, returnErr == nil)
	result, returnErr = api.tracersAPI.TraceBlockByHash(ctx, hash, config)
	return
}

func (api *DebugAPI) TraceCall(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, config *tracers.TraceCallConfig) (result interface{}, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("debug_traceCall", api.connectionType, startTime, returnErr == nil)
	result, returnErr = api.tracersAPI.TraceCall(ctx, args, blockNrOrHash, config)
	return
}
