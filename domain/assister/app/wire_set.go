package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/assister"
)

var New = wire.NewSet(
	NewAssisterController,
	NewAssisterRouter,
	domain.NewAssisterService,
)
