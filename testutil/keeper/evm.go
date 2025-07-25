package keeper

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bluelink-lab/blk-chain/app"
	evmkeeper "github.com/bluelink-lab/blk-chain/x/evm/keeper"
	evmtypes "github.com/bluelink-lab/blk-chain/x/evm/types"
)

var EVMTestApp = app.Setup(false, true)
var mockKeeper *evmkeeper.Keeper
var mockCtx sdk.Context
var mtx = &sync.Mutex{}

func MockEVMKeeperWithPrecompiles() (*evmkeeper.Keeper, sdk.Context) {
	mtx.Lock()
	defer mtx.Unlock()
	if mockKeeper != nil {
		return mockKeeper, mockCtx
	}
	ctx := EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockHeight(8)
	k := EVMTestApp.EvmKeeper
	k.InitGenesis(ctx, *evmtypes.DefaultGenesis())

	// mint some coins to a she address
	sheAddr, err := sdk.AccAddressFromHex(common.Bytes2Hex([]byte("sheAddr")))
	if err != nil {
		panic(err)
	}
	err = EVMTestApp.BankKeeper.MintCoins(ctx, "evm", sdk.NewCoins(sdk.NewCoin("ublt", sdk.NewInt(10))))
	if err != nil {
		panic(err)
	}
	err = EVMTestApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", sheAddr, sdk.NewCoins(sdk.NewCoin("ublt", sdk.NewInt(10))))
	if err != nil {
		panic(err)
	}
	mockKeeper = &k
	mockCtx = ctx
	return &k, ctx
}

func MockEVMKeeper() (*evmkeeper.Keeper, sdk.Context) {
	testApp := app.Setup(false, false)
	ctx := testApp.GetContextForDeliverTx([]byte{}).WithBlockHeight(8).WithBlockTime(time.Now())
	k := testApp.EvmKeeper
	k.InitGenesis(ctx, *evmtypes.DefaultGenesis())

	// mint some coins to a she address
	sheAddr, err := sdk.AccAddressFromHex(common.Bytes2Hex([]byte("sheAddr")))
	if err != nil {
		panic(err)
	}
	err = testApp.BankKeeper.MintCoins(ctx, "evm", sdk.NewCoins(sdk.NewCoin("ublt", sdk.NewInt(10))))
	if err != nil {
		panic(err)
	}
	err = testApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", sheAddr, sdk.NewCoins(sdk.NewCoin("ublt", sdk.NewInt(10))))
	if err != nil {
		panic(err)
	}
	return &k, ctx
}

func MockAddressPair() (sdk.AccAddress, common.Address) {
	return PrivateKeyToAddresses(MockPrivateKey())
}

func MockPrivateKey() cryptotypes.PrivKey {
	// Generate a new BLT private key
	entropySeed, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropySeed)
	algo := hd.Secp256k1
	derivedPriv, _ := algo.Derive()(mnemonic, "", "")
	return algo.Generate()(derivedPriv)
}

func PrivateKeyToAddresses(privKey cryptotypes.PrivKey) (sdk.AccAddress, common.Address) {
	// Encode the private key to hex (i.e. what wallets do behind the scene when users reveal private keys)
	testPrivHex := hex.EncodeToString(privKey.Bytes())

	// Sign an Ethereum transaction with the hex private key
	key, _ := crypto.HexToECDSA(testPrivHex)
	msg := crypto.Keccak256([]byte("foo"))
	sig, _ := crypto.Sign(msg, key)

	// Recover the public keys from the Ethereum signature
	recoveredPub, _ := crypto.Ecrecover(msg, sig)
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)

	return sdk.AccAddress(privKey.PubKey().Address()), crypto.PubkeyToAddress(*pubKey)
}

func UsheCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(sdk.MustGetBaseDenom(), sdk.NewInt(amount)))
}
