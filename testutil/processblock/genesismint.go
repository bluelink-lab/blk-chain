package processblock

import (
	"time"

	minttypes "github.com/bluelink-lab/blk-chain/x/mint/types"
)

func (a *App) NewMinter(amount uint64) {
	today := time.Now()
	dayAfterTomorrow := today.Add(48 * time.Hour)
	a.MintKeeper.SetMinter(a.Ctx(), minttypes.Minter{
		StartDate:           today.Format(minttypes.TokenReleaseDateFormat),
		EndDate:             dayAfterTomorrow.Format(minttypes.TokenReleaseDateFormat),
		Denom:               "ublt",
		TotalMintAmount:     amount,
		RemainingMintAmount: amount,
	})
}
