package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/mission"
)

var New = wire.NewSet(
	NewMissionController,
	NewMissionRouter,
	domain.NewMissionService,
)
