package state_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	testkeeper "github.com/bluelink-lab/blk-chain/testutil/keeper"
	"github.com/bluelink-lab/blk-chain/x/evm/state"
	"github.com/stretchr/testify/require"
)

func TestCode(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	_, addr := testkeeper.MockAddressPair()
	statedb := state.NewDBImpl(ctx, k, false)

	require.Equal(t, common.Hash{}, statedb.GetCodeHash(addr))
	require.Nil(t, statedb.GetCode(addr))
	require.Equal(t, 0, statedb.GetCodeSize(addr))

	code := []byte{1, 2, 3, 4, 5}
	statedb.SetCode(addr, code)
	require.Equal(t, crypto.Keccak256Hash(code), statedb.GetCodeHash(addr))
	require.Equal(t, code, statedb.GetCode(addr))
	require.Equal(t, 5, statedb.GetCodeSize(addr))
}
