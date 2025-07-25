package migrations

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bluelink-lab/blk-chain/x/evm/artifacts/erc1155"
	"github.com/bluelink-lab/blk-chain/x/evm/artifacts/erc20"
	"github.com/bluelink-lab/blk-chain/x/evm/artifacts/erc721"
	artifactsutils "github.com/bluelink-lab/blk-chain/x/evm/artifacts/utils"
	"github.com/bluelink-lab/blk-chain/x/evm/keeper"
	"github.com/bluelink-lab/blk-chain/x/evm/types"
)

func StoreCWPointerCode(ctx sdk.Context, k *keeper.Keeper, store20 bool, store721 bool, store1155 bool) error {
	if store20 {
		erc20CodeID, err := k.WasmKeeper().Create(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName), erc20.GetBin(), nil)
		if err != nil {
			panic(err)
		}
		prefix.NewStore(k.PrefixStore(ctx, types.PointerCWCodePrefix), types.PointerCW20ERC20Prefix).Set(
			artifactsutils.GetVersionBz(erc20.CurrentVersion),
			artifactsutils.GetCodeIDBz(erc20CodeID),
		)
	}

	if store721 {
		erc721CodeID, err := k.WasmKeeper().Create(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName), erc721.GetBin(), nil)
		if err != nil {
			panic(err)
		}
		prefix.NewStore(k.PrefixStore(ctx, types.PointerCWCodePrefix), types.PointerCW721ERC721Prefix).Set(
			artifactsutils.GetVersionBz(erc721.CurrentVersion),
			artifactsutils.GetCodeIDBz(erc721CodeID),
		)
	}

	if store1155 {
		erc1155CodeID, err := k.WasmKeeper().Create(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName), erc1155.GetBin(), nil)
		if err != nil {
			panic(err)
		}
		prefix.NewStore(k.PrefixStore(ctx, types.PointerCWCodePrefix), types.PointerCW1155ERC1155Prefix).Set(
			artifactsutils.GetVersionBz(erc1155.CurrentVersion),
			artifactsutils.GetCodeIDBz(erc1155CodeID),
		)
	}
	return nil
}
