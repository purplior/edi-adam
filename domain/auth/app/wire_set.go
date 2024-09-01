package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/auth"
)

var New = wire.NewSet(
	NewAuthRouter,
	NewAuthController,
	domain.NewAuthService,
)
