package state_test

import (
	"math/big"
	"testing"

	"github.com/bluelink-lab/blk-chain/x/evm/state"
	"github.com/stretchr/testify/require"
)

func TestGetCoinbaseAddress(t *testing.T) {
	coinbaseAddr := state.GetCoinbaseAddress(1).String()
	require.Equal(t, coinbaseAddr, "she1v4mx6hmrda5kucnpwdjsqqqqqqqqqqqpz6djs7")
}

func TestSplitUsheWeiAmount(t *testing.T) {
	for _, test := range []struct {
		amt         *big.Int
		expectedShe *big.Int
		expectedWei *big.Int
	}{
		{
			amt:         big.NewInt(0),
			expectedShe: big.NewInt(0),
			expectedWei: big.NewInt(0),
		}, {
			amt:         big.NewInt(1),
			expectedShe: big.NewInt(0),
			expectedWei: big.NewInt(1),
		}, {
			amt:         big.NewInt(999_999_999_999),
			expectedShe: big.NewInt(0),
			expectedWei: big.NewInt(999_999_999_999),
		}, {
			amt:         big.NewInt(1_000_000_000_000),
			expectedShe: big.NewInt(1),
			expectedWei: big.NewInt(0),
		}, {
			amt:         big.NewInt(1_000_000_000_001),
			expectedShe: big.NewInt(1),
			expectedWei: big.NewInt(1),
		}, {
			amt:         big.NewInt(123_456_789_123_456_789),
			expectedShe: big.NewInt(123456),
			expectedWei: big.NewInt(789_123_456_789),
		},
	} {
		ublk, wei := state.SplitUsheWeiAmount(test.amt)
		require.Equal(t, test.expectedShe, ublk.BigInt())
		require.Equal(t, test.expectedWei, wei.BigInt())
	}
}
