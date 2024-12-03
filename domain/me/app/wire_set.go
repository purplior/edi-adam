package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/me"
)

var New = wire.NewSet(
	NewMeController,
	NewMeRouter,
	domain.NewMeService,
)
