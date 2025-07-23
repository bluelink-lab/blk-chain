package main

import (
	"os"

	"github.com/bluelink-lab/blk-chain/app/params"
	"github.com/bluelink-lab/blk-chain/cmd/blkd/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/bluelink-lab/blk-chain/app"
)

func main() {
	params.SetAddressPrefixes()
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
