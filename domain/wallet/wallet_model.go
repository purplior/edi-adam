package wallet

import (
	"time"

	"github.com/podossaem/podoroot/domain/ledger"
)

type (
	Wallet struct {
		ID        string          `json:"id"`
		OwnerID   string          `json:"ownerId"`
		Podo      int             `json:"podo"`
		CreatedAt time.Time       `json:"createdAt"`
		Ledgers   []ledger.Ledger `json:"ledgers"`
	}
)
