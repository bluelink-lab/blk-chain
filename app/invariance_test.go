package app_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	app "github.com/bluelink-lab/blk-chain/app"
	"github.com/stretchr/testify/require"
)

func TestLightInvarianceChecks(t *testing.T) {
	tm := time.Now().UTC()
	valPub := secp256k1.GenPrivKey().PubKey()
	accounts := []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}
	ublkCoin := func(i int64) sdk.Coin { return sdk.NewCoin("ublk", sdk.NewInt(i)) }
	ublkCoins := func(i int64) sdk.Coins { return sdk.NewCoins(ublkCoin(i)) }
	for i, tt := range []struct {
		preUshe    []int64
		preWei     []int64
		preSupply  int64
		postUshe   []int64
		postWei    []int64
		postSupply int64
		success    bool
	}{
		{
			preUshe:    []int64{0, 0},
			preWei:     []int64{0, 0},
			preSupply:  5,
			postUshe:   []int64{1, 2},
			postWei:    []int64{0, 0},
			postSupply: 8,
			success:    true,
		},
		{
			preUshe:    []int64{2, 1},
			preWei:     []int64{0, 0},
			preSupply:  3,
			postUshe:   []int64{0, 0},
			postWei:    []int64{0, 0},
			postSupply: 0,
			success:    true,
		},
		{
			preUshe:    []int64{1, 0},
			preWei:     []int64{0, 0},
			preSupply:  10,
			postUshe:   []int64{0, 1},
			postWei:    []int64{0, 0},
			postSupply: 10,
			success:    true,
		},
		{
			preUshe:    []int64{1, 0},
			preWei:     []int64{0, 0},
			preSupply:  10,
			postUshe:   []int64{0, 0},
			postWei:    []int64{500_000_000_000, 500_000_000_000},
			postSupply: 10,
			success:    true,
		},
		{
			preUshe:    []int64{0, 0},
			preWei:     []int64{500_000_000_000, 500_000_000_000},
			preSupply:  10,
			postUshe:   []int64{1, 0},
			postWei:    []int64{0, 0},
			postSupply: 10,
			success:    true,
		},
		{
			preUshe:    []int64{0, 0},
			preWei:     []int64{1, 2},
			preSupply:  10,
			postUshe:   []int64{0, 0},
			postWei:    []int64{2, 1},
			postSupply: 10,
			success:    true,
		},
		{
			preUshe:    []int64{1, 0},
			preWei:     []int64{0, 0},
			preSupply:  10,
			postUshe:   []int64{1, 1},
			postWei:    []int64{0, 0},
			postSupply: 10,
			success:    false,
		},
		{
			preUshe:    []int64{1, 0},
			preWei:     []int64{0, 0},
			preSupply:  10,
			postUshe:   []int64{0, 0},
			postWei:    []int64{0, 0},
			postSupply: 10,
			success:    false,
		},
		{
			preUshe:    []int64{1, 0},
			preWei:     []int64{0, 0},
			preSupply:  10,
			postUshe:   []int64{0, 1},
			postWei:    []int64{500_000_000_000, 500_000_000_000},
			postSupply: 10,
			success:    false,
		},
		{
			preUshe:    []int64{1, 0},
			preWei:     []int64{500_000_000_000, 500_000_000_000},
			preSupply:  10,
			postUshe:   []int64{0, 1},
			postWei:    []int64{0, 0},
			postSupply: 10,
			success:    false,
		},
		{
			preUshe:    []int64{0, 0},
			preWei:     []int64{1, 2},
			preSupply:  10,
			postUshe:   []int64{0, 0},
			postWei:    []int64{2, 2},
			postSupply: 10,
			success:    false,
		},
		{
			preUshe:    []int64{0, 0},
			preWei:     []int64{1, 2},
			preSupply:  10,
			postUshe:   []int64{0, 0},
			postWei:    []int64{1, 1},
			postSupply: 10,
			success:    false,
		},
	} {
		fmt.Printf("Running test %d\n", i)
		testWrapper := app.NewTestWrapperWithSc(t, tm, valPub, false)
		a, ctx := testWrapper.App, testWrapper.Ctx
		for i := range tt.preUshe {
			if tt.preUshe[i] > 0 {
				a.BankKeeper.AddCoins(ctx, accounts[i], ublkCoins(tt.preUshe[i]), false)
			}
			if tt.preWei[i] > 0 {
				a.BankKeeper.AddWei(ctx, accounts[i], sdk.NewInt(tt.preWei[i]))
			}
		}
		if tt.preSupply > 0 {
			a.BankKeeper.SetSupply(ctx, ublkCoin(tt.preSupply))
		}
		a.SetDeliverStateToCommit()
		a.WriteState()
		a.GetWorkingHash() // flush to sc
		for i := range tt.postUshe {
			ublkDiff := tt.postUshe[i] - tt.preUshe[i]
			if ublkDiff > 0 {
				a.BankKeeper.AddCoins(ctx, accounts[i], ublkCoins(ublkDiff), false)
			} else if ublkDiff < 0 {
				a.BankKeeper.SubUnlockedCoins(ctx, accounts[i], ublkCoins(-ublkDiff), false)
			}

			weiDiff := tt.postWei[i] - tt.preWei[i]
			if weiDiff > 0 {
				a.BankKeeper.AddWei(ctx, accounts[i], sdk.NewInt(weiDiff))
			} else if weiDiff < 0 {
				a.BankKeeper.SubWei(ctx, accounts[i], sdk.NewInt(-weiDiff))
			}
		}
		a.BankKeeper.SetSupply(ctx, ublkCoin(tt.postSupply))
		a.SetDeliverStateToCommit()
		f := func() { a.LightInvarianceChecks(a.WriteState(), app.LightInvarianceConfig{SupplyEnabled: true}) }
		if tt.success {
			require.NotPanics(t, f)
		} else {
			require.Panics(t, f)
		}
		safeClose(a)
	}
}

// TODO: remove once snapshot manager can be closed gracefully in tests
func safeClose(a *app.App) {
	defer func() {
		_ = recover()
	}()
	a.Close()
}
