package keeper_test

import (
	"bytes"
	"testing"

	testkeeper "github.com/she-protocol/she-chain/testutil/keeper"
	"github.com/she-protocol/she-chain/x/evm/keeper"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{})
	// coinbase address must be associated
	coinbaseSheAddr, associated := k.GetSheAddress(ctx, keeper.GetCoinbaseAddress())
	require.True(t, associated)
	require.True(t, bytes.Equal(coinbaseSheAddr, k.AccountKeeper().GetModuleAddress("fee_collector")))
}
