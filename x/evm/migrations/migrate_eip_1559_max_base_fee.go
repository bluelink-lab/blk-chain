package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/she-protocol/she-chain/x/evm/keeper"
	"github.com/she-protocol/she-chain/x/evm/types"
)

func MigrateEip1559MaxFeePerGas(ctx sdk.Context, k *keeper.Keeper) error {
	keeperParams := k.GetParamsIfExists(ctx)
	keeperParams.MaximumFeePerGas = types.DefaultParams().MaximumFeePerGas
	k.SetParams(ctx, keeperParams)
	return nil
}
