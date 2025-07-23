package keeper_test

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/she-protocol/she-chain/testutil/keeper"
	evmkeeper "github.com/she-protocol/she-chain/x/evm/keeper"
	"github.com/stretchr/testify/require"
)

func TestSetGetAddressMapping(t *testing.T) {
	k := &keeper.EVMTestApp.EvmKeeper
	ctx := keeper.EVMTestApp.GetContextForDeliverTx([]byte{})
	sheAddr, evmAddr := keeper.MockAddressPair()
	_, ok := k.GetEVMAddress(ctx, sheAddr)
	require.False(t, ok)
	_, ok = k.GetSheAddress(ctx, evmAddr)
	require.False(t, ok)
	k.SetAddressMapping(ctx, sheAddr, evmAddr)
	foundEVM, ok := k.GetEVMAddress(ctx, sheAddr)
	require.True(t, ok)
	require.Equal(t, evmAddr, foundEVM)
	foundShe, ok := k.GetSheAddress(ctx, evmAddr)
	require.True(t, ok)
	require.Equal(t, sheAddr, foundShe)
	require.Equal(t, sheAddr, k.AccountKeeper().GetAccount(ctx, sheAddr).GetAddress())
}

func TestDeleteAddressMapping(t *testing.T) {
	k := &keeper.EVMTestApp.EvmKeeper
	ctx := keeper.EVMTestApp.GetContextForDeliverTx([]byte{})
	sheAddr, evmAddr := keeper.MockAddressPair()
	k.SetAddressMapping(ctx, sheAddr, evmAddr)
	foundEVM, ok := k.GetEVMAddress(ctx, sheAddr)
	require.True(t, ok)
	require.Equal(t, evmAddr, foundEVM)
	foundShe, ok := k.GetSheAddress(ctx, evmAddr)
	require.True(t, ok)
	require.Equal(t, sheAddr, foundShe)
	k.DeleteAddressMapping(ctx, sheAddr, evmAddr)
	_, ok = k.GetEVMAddress(ctx, sheAddr)
	require.False(t, ok)
	_, ok = k.GetSheAddress(ctx, evmAddr)
	require.False(t, ok)
}

func TestGetAddressOrDefault(t *testing.T) {
	k := &keeper.EVMTestApp.EvmKeeper
	ctx := keeper.EVMTestApp.GetContextForDeliverTx([]byte{})
	sheAddr, evmAddr := keeper.MockAddressPair()
	defaultEvmAddr := k.GetEVMAddressOrDefault(ctx, sheAddr)
	require.True(t, bytes.Equal(sheAddr, defaultEvmAddr[:]))
	defaultSheAddr := k.GetSheAddressOrDefault(ctx, evmAddr)
	require.True(t, bytes.Equal(defaultSheAddr, evmAddr[:]))
}

func TestSendingToCastAddress(t *testing.T) {
	a := keeper.EVMTestApp
	ctx := a.GetContextForDeliverTx([]byte{})
	sheAddr, evmAddr := keeper.MockAddressPair()
	castAddr := sdk.AccAddress(evmAddr[:])
	sourceAddr, _ := keeper.MockAddressPair()
	require.Nil(t, a.BankKeeper.MintCoins(ctx, "evm", sdk.NewCoins(sdk.NewCoin("ublk", sdk.NewInt(10)))))
	require.Nil(t, a.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", sourceAddr, sdk.NewCoins(sdk.NewCoin("ublk", sdk.NewInt(5)))))
	amt := sdk.NewCoins(sdk.NewCoin("ublk", sdk.NewInt(1)))
	require.Nil(t, a.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", castAddr, amt))
	require.Nil(t, a.BankKeeper.SendCoins(ctx, sourceAddr, castAddr, amt))
	require.Nil(t, a.BankKeeper.SendCoinsAndWei(ctx, sourceAddr, castAddr, sdk.OneInt(), sdk.OneInt()))

	a.EvmKeeper.SetAddressMapping(ctx, sheAddr, evmAddr)
	require.NotNil(t, a.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", castAddr, amt))
	require.NotNil(t, a.BankKeeper.SendCoins(ctx, sourceAddr, castAddr, amt))
	require.NotNil(t, a.BankKeeper.SendCoinsAndWei(ctx, sourceAddr, castAddr, sdk.OneInt(), sdk.OneInt()))
}

func TestEvmAddressHandler_GetSheAddressFromString(t *testing.T) {
	a := keeper.EVMTestApp
	ctx := a.GetContextForDeliverTx([]byte{})
	sheAddr, evmAddr := keeper.MockAddressPair()
	a.EvmKeeper.SetAddressMapping(ctx, sheAddr, evmAddr)

	_, notAssociatedEvmAddr := keeper.MockAddressPair()
	castAddr := sdk.AccAddress(notAssociatedEvmAddr[:])

	type args struct {
		ctx     sdk.Context
		address string
	}
	tests := []struct {
		name       string
		args       args
		want       sdk.AccAddress
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "returns associated BLK address if input address is a valid 0x and associated",
			args: args{
				ctx:     ctx,
				address: evmAddr.String(),
			},
			want: sheAddr,
		},
		{
			name: "returns default BLK address if input address is a valid 0x not associated",
			args: args{
				ctx:     ctx,
				address: notAssociatedEvmAddr.String(),
			},
			want: castAddr,
		},
		{
			name: "returns BLK address if input address is a valid bech32 address",
			args: args{
				ctx:     ctx,
				address: sheAddr.String(),
			},
			want: sheAddr,
		},
		{
			name: "returns error if address is invalid",
			args: args{
				ctx:     ctx,
				address: "invalid",
			},
			wantErr:    true,
			wantErrMsg: "decoding bech32 failed: invalid bech32 string length 7",
		}, {
			name: "returns error if address is empty",
			args: args{
				ctx:     ctx,
				address: "",
			},
			wantErr:    true,
			wantErrMsg: "empty address string is not allowed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := evmkeeper.NewEvmAddressHandler(&a.EvmKeeper)
			got, err := h.GetSheAddressFromString(tt.args.ctx, tt.args.address)
			if tt.wantErr {
				require.NotNil(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
