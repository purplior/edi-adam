package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/mission"
)

var New = wire.NewSet(
	NewMissionController,
	NewMissionRouter,
	domain.NewMissionService,
)
