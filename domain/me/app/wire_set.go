package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/me"
)

var New = wire.NewSet(
	NewMeController,
	NewMeRouter,
	domain.NewMeService,
)
