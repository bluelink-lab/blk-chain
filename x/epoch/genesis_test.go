package epoch_test

import (
	"testing"
	"time"

	keepertest "github.com/bluelink-lab/blk-chain/testutil/keeper"
	"github.com/bluelink-lab/blk-chain/testutil/nullify"
	"github.com/bluelink-lab/blk-chain/x/epoch"
	"github.com/bluelink-lab/blk-chain/x/epoch/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	now := time.Now()
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Epoch: &types.Epoch{
			GenesisTime:           now,
			EpochDuration:         time.Minute,
			CurrentEpoch:          1,
			CurrentEpochStartTime: now,
			CurrentEpochHeight:    0,
		},
	}

	k, ctx := keepertest.EpochKeeper(t)
	epoch.InitGenesis(ctx, *k, genesisState)
	got := epoch.ExportGenesis(ctx, *k)
	require.NotNil(t, got)
	require.Equal(t, got.Epoch.CurrentEpoch, genesisState.Epoch.CurrentEpoch)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
