package migrations_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/bluelink-lab/blk-chain/testutil/keeper"
	"github.com/bluelink-lab/blk-chain/x/evm/migrations"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMigrateEip1559MaxBaseFee(t *testing.T) {
	k := testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.NewContext(false, tmtypes.Header{})

	keeperParams := k.GetParams(ctx)
	keeperParams.MaximumFeePerGas = sdk.NewDec(123)

	// Perform the migration
	err := migrations.MigrateEip1559MaxFeePerGas(ctx, &k)
	require.NoError(t, err)

	// Ensure that the new EIP-1559 parameters were migrated and the old ones were not changed
	require.Equal(t, keeperParams.MaximumFeePerGas, sdk.NewDec(123))
}
