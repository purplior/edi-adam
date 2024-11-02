package entity

import (
	"time"

	"github.com/podossaem/podoroot/domain/ledger"
	domain "github.com/podossaem/podoroot/domain/wallet"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Wallet struct {
		ID        uint      `gorm:"primaryKey;autoIncrement"`
		OwnerID   uint      `gorm:"unique"`
		Podo      int64     `gorm:"not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		Ledgers   []Ledger  `gorm:"foreignKey:WalletID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e Wallet) ToModel() domain.Wallet {
	m := domain.Wallet{
		Podo:      dt.Int(e.Podo),
		CreatedAt: e.CreatedAt,
	}

	if e.ID > 0 {
		m.ID = dt.Str(e.ID)
	}
	if e.OwnerID > 0 {
		m.OwnerID = dt.Str(e.OwnerID)
	}

	m.Ledgers = make([]ledger.Ledger, len(e.Ledgers))
	for i, ledger := range e.Ledgers {
		m.Ledgers[i] = ledger.ToModel()
	}

	return m
}

func MakeWallet(m domain.Wallet) Wallet {
	e := Wallet{
		Podo:      dt.Int64(m.Podo),
		CreatedAt: m.CreatedAt,
	}

	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}
	if len(m.OwnerID) > 0 {
		e.OwnerID = dt.UInt(m.OwnerID)
	}

	e.Ledgers = make([]Ledger, len(m.Ledgers))
	for i, ledger := range m.Ledgers {
		e.Ledgers[i] = MakeLedger(ledger)
	}

	return e
}
