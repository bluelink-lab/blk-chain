package migrations_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/bluelink-lab/blk-chain/testutil/keeper"
	"github.com/bluelink-lab/blk-chain/x/evm/migrations"
	"github.com/bluelink-lab/blk-chain/x/evm/types"
	"github.com/stretchr/testify/require"
)

func TestMigrateCastAddressBalances(t *testing.T) {
	k := testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	require.Nil(t, k.BankKeeper().MintCoins(ctx, types.ModuleName, testkeeper.UsheCoins(100)))
	// unassociated account with funds
	sheAddr1, evmAddr1 := testkeeper.MockAddressPair()
	require.Nil(t, k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(evmAddr1[:]), testkeeper.UsheCoins(10)))
	// associated account without funds
	sheAddr2, evmAddr2 := testkeeper.MockAddressPair()
	k.SetAddressMapping(ctx, sheAddr2, evmAddr2)
	// associated account with funds
	sheAddr3, evmAddr3 := testkeeper.MockAddressPair()
	require.Nil(t, k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(evmAddr3[:]), testkeeper.UsheCoins(10)))
	k.SetAddressMapping(ctx, sheAddr3, evmAddr3)

	require.Nil(t, migrations.MigrateCastAddressBalances(ctx, &k))

	require.Equal(t, sdk.NewInt(10), k.BankKeeper().GetBalance(ctx, sdk.AccAddress(evmAddr1[:]), "ublk").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, sheAddr1, "ublk").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, sdk.AccAddress(evmAddr2[:]), "ublk").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, sheAddr2, "ublk").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, sdk.AccAddress(evmAddr3[:]), "ublk").Amount)
	require.Equal(t, sdk.NewInt(10), k.BankKeeper().GetBalance(ctx, sheAddr3, "ublk").Amount)
}
