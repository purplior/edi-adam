package app

import (
	"github.com/google/wire"
	"github.com/purplior/sbec/domain/user"
)

var New = wire.NewSet(
	NewUserRouter,
	NewUserController,
	user.NewUserService,
)
