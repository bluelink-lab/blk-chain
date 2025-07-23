package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bluelink-lab/blk-chain/x/evm/state"
)

func (k *Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) *big.Int {
	denom := k.GetBaseDenom(ctx)
	allUshe := k.BankKeeper().GetBalance(ctx, addr, denom).Amount
	lockedUshe := k.BankKeeper().LockedCoins(ctx, addr).AmountOf(denom) // LockedCoins doesn't use iterators
	ublt := allUshe.Sub(lockedUshe)
	wei := k.BankKeeper().GetWeiBalance(ctx, addr)
	return ublt.Mul(state.SdkUsheToSweiMultiplier).Add(wei).BigInt()
}
