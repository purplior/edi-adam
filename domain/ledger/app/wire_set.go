package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/ledger"
)

var New = wire.NewSet(
	domain.NewLedgerService,
)
