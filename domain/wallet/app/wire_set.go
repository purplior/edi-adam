package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/wallet"
)

var New = wire.NewSet(
	domain.NewWalletService,
)
