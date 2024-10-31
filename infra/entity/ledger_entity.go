package entity

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/ledger"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Ledger struct {
		ID         uint `gorm:"primaryKey;autoIncrement"`
		WalletID   uint
		PodoAmount int32 `gorm:"not null"`
		Action     domain.LedgerAction
		Reason     string    `gorm:"size:255"`
		CreatedAt  time.Time `gorm:"autoCreateTime"`
	}
)

func (e Ledger) ToModel() domain.Ledger {
	m := domain.Ledger{
		PodoAmount: dt.Int(e.PodoAmount),
		Action:     e.Action,
		Reason:     e.Reason,
		CreatedAt:  e.CreatedAt,
	}

	if e.ID > 0 {
		m.ID = dt.Str(e.ID)
	}
	if e.WalletID > 0 {
		m.WalletID = dt.Str(e.WalletID)
	}

	return m
}

func MakeLedger(m domain.Ledger) Ledger {
	e := Ledger{
		PodoAmount: dt.Int32(m.PodoAmount),
		Action:     m.Action,
		Reason:     m.Reason,
		CreatedAt:  m.CreatedAt,
	}

	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}
	if len(m.WalletID) > 0 {
		e.WalletID = dt.UInt(m.WalletID)
	}

	return e
}
