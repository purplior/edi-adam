package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/assisterform"
)

var New = wire.NewSet(
	NewAssisterFormController,
	NewAssisterFormRouter,
	domain.NewAssisterFormService,
)
