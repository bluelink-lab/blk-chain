package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/she-protocol/she-chain/x/evm/types"
)

func (k *Keeper) SetAddressMapping(ctx sdk.Context, sheAddress sdk.AccAddress, evmAddress common.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.EVMAddressToSheAddressKey(evmAddress), sheAddress)
	store.Set(types.SheAddressToEVMAddressKey(sheAddress), evmAddress[:])
	if !k.accountKeeper.HasAccount(ctx, sheAddress) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, sheAddress))
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAddressAssociated,
		sdk.NewAttribute(types.AttributeKeySheAddress, sheAddress.String()),
		sdk.NewAttribute(types.AttributeKeyEvmAddress, evmAddress.Hex()),
	))
}

func (k *Keeper) DeleteAddressMapping(ctx sdk.Context, sheAddress sdk.AccAddress, evmAddress common.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.EVMAddressToSheAddressKey(evmAddress))
	store.Delete(types.SheAddressToEVMAddressKey(sheAddress))
}

func (k *Keeper) GetEVMAddress(ctx sdk.Context, sheAddress sdk.AccAddress) (common.Address, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SheAddressToEVMAddressKey(sheAddress))
	addr := common.Address{}
	if bz == nil {
		return addr, false
	}
	copy(addr[:], bz)
	return addr, true
}

func (k *Keeper) GetEVMAddressOrDefault(ctx sdk.Context, sheAddress sdk.AccAddress) common.Address {
	addr, ok := k.GetEVMAddress(ctx, sheAddress)
	if ok {
		return addr
	}
	return common.BytesToAddress(sheAddress)
}

func (k *Keeper) GetSheAddress(ctx sdk.Context, evmAddress common.Address) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EVMAddressToSheAddressKey(evmAddress))
	if bz == nil {
		return []byte{}, false
	}
	return bz, true
}

func (k *Keeper) GetSheAddressOrDefault(ctx sdk.Context, evmAddress common.Address) sdk.AccAddress {
	addr, ok := k.GetSheAddress(ctx, evmAddress)
	if ok {
		return addr
	}
	return sdk.AccAddress(evmAddress[:])
}

func (k *Keeper) IterateSheAddressMapping(ctx sdk.Context, cb func(evmAddr common.Address, sheAddr sdk.AccAddress) bool) {
	iter := prefix.NewStore(ctx.KVStore(k.storeKey), types.EVMAddressToSheAddressKeyPrefix).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		evmAddr := common.BytesToAddress(iter.Key())
		sheAddr := sdk.AccAddress(iter.Value())
		if cb(evmAddr, sheAddr) {
			break
		}
	}
}

// A sdk.AccAddress may not receive funds from bank if it's the result of direct-casting
// from an EVM address AND the originating EVM address has already been associated with
// a true (i.e. derived from the same pubkey) sdk.AccAddress.
func (k *Keeper) CanAddressReceive(ctx sdk.Context, addr sdk.AccAddress) bool {
	directCast := common.BytesToAddress(addr) // casting goes both directions since both address formats have 20 bytes
	associatedAddr, isAssociated := k.GetSheAddress(ctx, directCast)
	// if the associated address is the cast address itself, allow the address to receive (e.g. EVM contract addresses)
	return associatedAddr.Equals(addr) || !isAssociated // this means it's either a cast address that's not associated yet, or not a cast address at all.
}

type EvmAddressHandler struct {
	evmKeeper *Keeper
}

func NewEvmAddressHandler(evmKeeper *Keeper) EvmAddressHandler {
	return EvmAddressHandler{evmKeeper: evmKeeper}
}

func (h EvmAddressHandler) GetSheAddressFromString(ctx sdk.Context, address string) (sdk.AccAddress, error) {
	if common.IsHexAddress(address) {
		parsedAddress := common.HexToAddress(address)
		return h.evmKeeper.GetSheAddressOrDefault(ctx, parsedAddress), nil
	}
	parsedAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	return parsedAddress, nil
}
