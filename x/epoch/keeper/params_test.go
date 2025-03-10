package keeper_test

import (
	"testing"

	testkeeper "github.com/she-protocol/she-chain/testutil/keeper"
	"github.com/she-protocol/she-chain/x/epoch/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.EpochKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
