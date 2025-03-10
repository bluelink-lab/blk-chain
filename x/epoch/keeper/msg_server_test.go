package keeper_test

import (
	"testing"

	keepertest "github.com/she-protocol/she-chain/testutil/keeper"
	"github.com/she-protocol/she-chain/x/epoch/keeper"
	"github.com/stretchr/testify/require"
)

func TestSetupMsgServer(t *testing.T) {
	k, _ := keepertest.EpochKeeper(t)
	msg := keeper.NewMsgServerImpl(*k)
	require.NotNil(t, msg)
}
