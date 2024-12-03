package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/assister"
)

var New = wire.NewSet(
	NewAssisterController,
	NewAssisterRouter,
	domain.NewAssisterService,
)
