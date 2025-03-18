package processblock

import (
	"time"

	minttypes "github.com/she-protocol/she-chain/x/mint/types"
)

func (a *App) NewMinter(amount uint64) {
	today := time.Now()
	dayAfterTomorrow := today.Add(48 * time.Hour)
	a.MintKeeper.SetMinter(a.Ctx(), minttypes.Minter{
		StartDate:           today.Format(minttypes.TokenReleaseDateFormat),
		EndDate:             dayAfterTomorrow.Format(minttypes.TokenReleaseDateFormat),
		Denom:               "ushe",
		TotalMintAmount:     amount,
		RemainingMintAmount: amount,
	})
}
