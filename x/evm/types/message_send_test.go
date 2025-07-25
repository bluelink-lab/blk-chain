package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/bluelink-lab/blk-chain/x/evm/types"
	"github.com/stretchr/testify/require"
)

func TestMessageSendValidate(t *testing.T) {
	fromAddr, err := sdk.AccAddressFromBech32("she1yezq49upxhunjjhudql2fnj5dgvcwjj87pn2wx")
	require.Nil(t, err)
	msg := types.NewMsgSend(fromAddr, common.HexToAddress("to"), sdk.Coins{sdk.Coin{
		Denom:  "blt",
		Amount: sdk.NewInt(1),
	}})
	require.Nil(t, msg.ValidateBasic())

	// No coins
	msg = types.NewMsgSend(fromAddr, common.HexToAddress("to"), sdk.Coins{})
	require.Error(t, msg.ValidateBasic())

	// Negative coins
	msg = types.NewMsgSend(fromAddr, common.HexToAddress("to"), sdk.Coins{sdk.Coin{
		Denom:  "blt",
		Amount: sdk.NewInt(-1),
	}})
	require.Error(t, msg.ValidateBasic())
}
