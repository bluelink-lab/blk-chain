package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/she-protocol/she-chain/x/evm/keeper"
)

func MigrateBaseFeeOffByOne(ctx sdk.Context, k *keeper.Keeper) error {
	baseFee := k.GetCurrBaseFeePerGas(ctx)
	k.SetNextBaseFeePerGas(ctx, baseFee)
	return nil
}
