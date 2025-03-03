package controller

import "github.com/google/wire"

var New = wire.NewSet(
	NewAssistantAdminController,
)
