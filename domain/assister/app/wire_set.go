package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/assister"
)

var New = wire.NewSet(
	NewAssisterController,
	NewAssisterRouter,
	domain.NewAssisterService,
)
