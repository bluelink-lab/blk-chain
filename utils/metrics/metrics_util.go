package metrics

import (
	"errors"
	"math/big"
	"strconv"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/bluelink-lab/blk-chain/x/evm/types"
)

// Measures the time taken to execute a sudo msg
// Metric Names:
//
//	blt_sudo_duration_miliseconds
//	blt_sudo_duration_miliseconds_count
//	blt_sudo_duration_miliseconds_sum
func MeasureSudoExecutionDuration(start time.Time, msgType string) {
	metrics.MeasureSinceWithLabels(
		[]string{"blt", "sudo", "duration", "milliseconds"},
		start.UTC(),
		[]metrics.Label{telemetry.NewLabel("type", msgType)},
	)
}

// Measures failed sudo execution count
// Metric Name:
//
//	blt_sudo_error_count
func IncrementSudoFailCount(msgType string) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "sudo", "error", "count"},
		1,
		[]metrics.Label{telemetry.NewLabel("type", msgType)},
	)
}

// Gauge metric with blkd version and git commit as labels
// Metric Name:
//
//	shed_version_and_commit
func GaugeShedVersionAndCommit(version string, commit string) {
	telemetry.SetGaugeWithLabels(
		[]string{"shed_version_and_commit"},
		1,
		[]metrics.Label{telemetry.NewLabel("shed_version", version), telemetry.NewLabel("commit", commit)},
	)
}

// blt_tx_process_type_count
func IncrTxProcessTypeCounter(processType string) {
	metrics.IncrCounterWithLabels(
		[]string{"blt", "tx", "process", "type"},
		1,
		[]metrics.Label{telemetry.NewLabel("type", processType)},
	)
}

// Measures the time taken to process a block by the process type
// Metric Names:
//
//	blt_process_block_miliseconds
//	blt_process_block_miliseconds_count
//	blt_process_block_miliseconds_sum
func BlockProcessLatency(start time.Time, processType string) {
	metrics.MeasureSinceWithLabels(
		[]string{"blt", "process", "block", "milliseconds"},
		start.UTC(),
		[]metrics.Label{telemetry.NewLabel("type", processType)},
	)
}

// Measures the time taken to execute a sudo msg
// Metric Names:
//
//	blt_tx_process_type_count
func IncrDagBuildErrorCounter(reason string) {
	metrics.IncrCounterWithLabels(
		[]string{"blt", "dag", "build", "error"},
		1,
		[]metrics.Label{telemetry.NewLabel("reason", reason)},
	)
}

// Counts the number of concurrent transactions that failed
// Metric Names:
//
//	blt_tx_concurrent_delivertx_error
func IncrFailedConcurrentDeliverTxCounter() {
	metrics.IncrCounterWithLabels(
		[]string{"blt", "tx", "concurrent", "delievertx", "error"},
		1,
		[]metrics.Label{},
	)
}

// Counts the number of operations that failed due to operation timeout
// Metric Names:
//
//	blt_log_not_done_after_counter
func IncrLogIfNotDoneAfter(label string) {
	metrics.IncrCounterWithLabels(
		[]string{"blt", "log", "not", "done", "after"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("label", label),
		},
	)
}

// Measures the time taken to execute a sudo msg
// Metric Names:
//
//	blt_deliver_tx_duration_miliseconds
//	blt_deliver_tx_duration_miliseconds_count
//	blt_deliver_tx_duration_miliseconds_sum
func MeasureDeliverTxDuration(start time.Time) {
	metrics.MeasureSince(
		[]string{"blt", "deliver", "tx", "milliseconds"},
		start.UTC(),
	)
}

// Measures the time taken to execute a batch tx
// Metric Names:
//
//	blt_deliver_batch_tx_duration_miliseconds
//	blt_deliver_batch_tx_duration_miliseconds_count
//	blt_deliver_batch_tx_duration_miliseconds_sum
func MeasureDeliverBatchTxDuration(start time.Time) {
	metrics.MeasureSince(
		[]string{"blt", "deliver", "batch", "tx", "milliseconds"},
		start.UTC(),
	)
}

// blt_oracle_vote_penalty_count
func SetOracleVotePenaltyCount(count uint64, valAddr string, penaltyType string) {
	metrics.SetGaugeWithLabels(
		[]string{"blt", "oracle", "vote", "penalty", "count"},
		float32(count),
		[]metrics.Label{
			telemetry.NewLabel("type", penaltyType),
			telemetry.NewLabel("validator", valAddr),
		},
	)
}

// blt_epoch_new
func SetEpochNew(epochNum uint64) {
	metrics.SetGauge(
		[]string{"blt", "epoch", "new"},
		float32(epochNum),
	)
}

// Measures throughput
// Metric Name:
//
//	blt_throughput_<metric_name>
func SetThroughputMetric(metricName string, value float32) {
	telemetry.SetGauge(
		value,
		"blt", "throughput", metricName,
	)
}

// Measures number of new websocket connects
// Metric Name:
//
//	blt_websocket_connect
func IncWebsocketConnects() {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "websocket", "connect"},
		1,
		nil,
	)
}

// Measures number of times a denom's price is updated
// Metric Name:
//
//	blt_oracle_price_update_count
func IncrPriceUpdateDenom(denom string) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "oracle", "price", "update"},
		1,
		[]metrics.Label{telemetry.NewLabel("denom", denom)},
	)
}

// Measures throughput per message type
// Metric Name:
//
//	blt_throughput_<metric_name>
func SetThroughputMetricByType(metricName string, value float32, msgType string) {
	telemetry.SetGaugeWithLabels(
		[]string{"blt", "loadtest", "tps", metricName},
		value,
		[]metrics.Label{telemetry.NewLabel("msg_type", msgType)},
	)
}

// Measures the number of times the total block gas wanted in the proposal exceeds the max
// Metric Name:
//
//	blt_failed_total_gas_wanted_check
func IncrFailedTotalGasWantedCheck(proposer string) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "failed", "total", "gas", "wanted", "check"},
		1,
		[]metrics.Label{telemetry.NewLabel("proposer", proposer)},
	)
}

// Measures the number of times the total block gas wanted in the proposal exceeds the max
// Metric Name:
//
//	blt_failed_total_gas_wanted_check
func IncrValidatorSlashed(proposer string) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "failed", "total", "gas", "wanted", "check"},
		1,
		[]metrics.Label{telemetry.NewLabel("proposer", proposer)},
	)
}

// Measures number of times a denom's price is updated
// Metric Name:
//
//	blt_oracle_price_update_count
func SetCoinsMinted(amount uint64, denom string) {
	telemetry.SetGaugeWithLabels(
		[]string{"blt", "mint", "coins"},
		float32(amount),
		[]metrics.Label{telemetry.NewLabel("denom", denom)},
	)
}

// Measures the number of times the total block gas wanted in the proposal exceeds the max
// Metric Name:
//
//	blt_tx_gas_counter
func IncrGasCounter(gasType string, value int64) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "tx", "gas", "counter"},
		float32(value),
		[]metrics.Label{telemetry.NewLabel("type", gasType)},
	)
}

// Measures the number of times optimistic processing runs
// Metric Name:
//
//	blt_optimistic_processing_counter
func IncrementOptimisticProcessingCounter(enabled bool) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "optimistic", "processing", "counter"},
		float32(1),
		[]metrics.Label{telemetry.NewLabel("enabled", strconv.FormatBool(enabled))},
	)
}

// Measures RPC endpoint request throughput
// Metric Name:
//
//	blt_rpc_request_counter
func IncrementRpcRequestCounter(endpoint string, connectionType string, success bool) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "rpc", "request", "counter"},
		float32(1),
		[]metrics.Label{
			telemetry.NewLabel("endpoint", endpoint),
			telemetry.NewLabel("connection", connectionType),
			telemetry.NewLabel("success", strconv.FormatBool(success)),
		},
	)
}

func IncrementErrorMetrics(scenario string, err error) {
	if err == nil {
		return
	}
	var assocErr types.AssociationMissingErr
	if errors.As(err, &assocErr) {
		IncrementAssociationError(scenario, assocErr)
		return
	}
	// add other error types to handle as metrics
}

func IncrementAssociationError(scenario string, err types.AssociationMissingErr) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "association", "error"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("scenario", scenario),
			telemetry.NewLabel("type", err.AddressType()),
		},
	)
}

// Measures the RPC request latency in milliseconds
// Metric Name:
//
//	blt_rpc_request_latency_ms
func MeasureRpcRequestLatency(endpoint string, connectionType string, startTime time.Time) {
	metrics.MeasureSinceWithLabels(
		[]string{"blt", "rpc", "request", "latency_ms"},
		startTime.UTC(),
		[]metrics.Label{
			telemetry.NewLabel("endpoint", endpoint),
			telemetry.NewLabel("connection", connectionType),
		},
	)
}

// IncrProducerEventCount increments the counter for events produced.
// This metric counts the number of events produced by the system.
// Metric Name:
//
//	blt_loadtest_produce_count
func IncrProducerEventCount(msgType string) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "loadtest", "produce", "count"},
		1,
		[]metrics.Label{telemetry.NewLabel("msg_type", msgType)},
	)
}

// IncrConsumerEventCount increments the counter for events consumed.
// This metric counts the number of events consumed by the system.
// Metric Name:
//
//	blt_loadtest_consume_count
func IncrConsumerEventCount(msgType string) {
	telemetry.IncrCounterWithLabels(
		[]string{"blt", "loadtest", "consume", "count"},
		1,
		[]metrics.Label{telemetry.NewLabel("msg_type", msgType)},
	)
}

func AddHistogramMetric(key []string, value float32) {
	metrics.AddSample(key, value)
}

// Gauge for gas price paid for transactions
// Metric Name:
//
// blt_evm_effective_gas_price
func HistogramEvmEffectiveGasPrice(gasPrice *big.Int) {
	AddHistogramMetric(
		[]string{"blt", "evm", "effective", "gas", "price"},
		float32(gasPrice.Uint64()),
	)
}

// Gauge for block base fee
// Metric Name:
//
// blt_evm_block_base_fee
func GaugeEvmBlockBaseFee(baseFee *big.Int, blockHeight int64) {
	metrics.SetGauge(
		[]string{"blt", "evm", "block", "base", "fee"},
		float32(baseFee.Uint64()),
	)
}
