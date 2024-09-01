package app

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/domain/user/persist"
)

var New = wire.NewSet(
	NewUserRouter,
	NewUserController,
	persist.NewUserRepository,
	user.NewUserService,
)
