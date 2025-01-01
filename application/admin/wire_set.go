package admin

import "github.com/google/wire"

var New = wire.NewSet(
	NewAdminRouter,
	NewAdminController,
)
