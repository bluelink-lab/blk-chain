package main

import (
	"os"

	"github.com/she-protocol/she-chain/app/params"
	"github.com/she-protocol/she-chain/cmd/blkd/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/she-protocol/she-chain/app"
)

func main() {
	params.SetAddressPrefixes()
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
