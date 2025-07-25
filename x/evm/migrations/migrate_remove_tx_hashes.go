package migrations

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bluelink-lab/blk-chain/x/evm/keeper"
	"github.com/bluelink-lab/blk-chain/x/evm/types"
)

func RemoveTxHashes(ctx sdk.Context, k *keeper.Keeper) error {
	store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), types.TxHashesPrefix)
	return store.DeleteAll(nil, nil)
}
