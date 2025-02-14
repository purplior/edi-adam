package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/auth"
)

var New = wire.NewSet(
	NewAuthRouter,
	NewAuthController,
	domain.NewAuthService,
)
