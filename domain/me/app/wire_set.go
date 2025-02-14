package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/me"
)

var New = wire.NewSet(
	NewMeController,
	NewMeRouter,
	domain.NewMeService,
)
