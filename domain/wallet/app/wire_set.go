package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/wallet"
)

var New = wire.NewSet(
	domain.NewWalletService,
)
