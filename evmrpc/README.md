# BLT's EVM RPC

BLT supports the standard [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/) endpoints. On top of that, BLT supports some additional custom endpoints.

## She_ endpoints

### Endpoints for Synthetic txs
The motivation for these endpoints is to expose CW20 and CW721 events on the EVM side through synthetic receipts and logs. This is useful for indexing pointer contracts.
 - `blt_getFilterLogs`
   - same as `eth_getFilterLogs` but includes synthetic logs
 - `blt_getLogs`
   - same as `eth_getLogs` but includes synthetic logs
 - `blt_getBlockByNumber` and `blt_getBlockByHash`
   - same as `eth_getBlockByNumber` and `eth_getBlockByHash` but includes synthetic txs
 - NOTE: for synthetic txs, `eth_getTransactionReceipt` can be used to get the receipt data for a synthetic tx hash.

### Endpoints for excluding tracing failures
The motivation for these endpoints is to exclude tracing failures from the EVM side. Due to how our mempool works and our lack of tx simulation, we cannot rely on txs to pass all pre-state checks. Therefore, in the eth_ endpoints, we may see txs that fail tracing with errors like "nonce too low", "nonce too high", "insufficient funds", or other types of panic failures. These transactions are not executed, yet are still included in the block. These endpoints are useful for filtering out these txs.
- `blt_traceBlockByNumberExcludeTraceFail`
  - same as `debug_traceBlockByNumber` but excludes panic txs
- `blt_traceBlockByHashExcludeTraceFail`
  - same as `debug_traceBlockByHash` but excludes panic txs
- `blt_getTransactionReceiptExcludeTraceFail`
  - same as `eth_getTransactionReceipt` but excludes panic txs
- `blt_getBlockByNumberExcludeTraceFail` and `blt_getBlockByHashExcludeTraceFail`
  - same as `eth_getBlockByNumber` and `eth_getBlockByHash` but excludes panic txs
