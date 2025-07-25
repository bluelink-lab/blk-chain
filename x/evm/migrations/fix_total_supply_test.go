package migrations_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/bluelink-lab/blk-chain/testutil/keeper"
	"github.com/bluelink-lab/blk-chain/x/evm/migrations"
	"github.com/stretchr/testify/require"
)

func TestFixTotalSupply(t *testing.T) {
	k := &testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	addr, _ := testkeeper.MockAddressPair()
	balance := sdk.NewCoins(sdk.NewCoin(sdk.MustGetBaseDenom(), sdk.OneInt()))
	k.BankKeeper().MintCoins(ctx, "evm", balance)
	k.BankKeeper().SendCoinsFromModuleToAccount(ctx, "evm", addr, balance)
	k.BankKeeper().AddWei(ctx, addr, sdk.OneInt())
	oldSupply := k.BankKeeper().GetSupply(ctx, sdk.MustGetBaseDenom()).Amount
	require.Nil(t, migrations.FixTotalSupply(ctx, k))
	require.Equal(t, oldSupply.Add(sdk.OneInt()), k.BankKeeper().GetSupply(ctx, sdk.MustGetBaseDenom()).Amount)
	require.Equal(t, sdk.OneInt(), k.BankKeeper().GetBalance(ctx, addr, sdk.MustGetBaseDenom()).Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, k.AccountKeeper().GetModuleAddress("evm"), sdk.MustGetBaseDenom()).Amount)
	require.Equal(t, sdk.OneInt(), k.BankKeeper().GetWeiBalance(ctx, addr))
	require.Equal(t, sdk.NewInt(999_999_999_999), k.BankKeeper().GetWeiBalance(ctx, k.AccountKeeper().GetModuleAddress("evm")))
}
