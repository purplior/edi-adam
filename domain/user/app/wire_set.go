package app

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/domain/user"
)

var New = wire.NewSet(
	NewUserRouter,
	NewUserController,
	user.NewUserService,
)
