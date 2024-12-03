package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/assisterform"
)

var New = wire.NewSet(
	NewAssisterFormController,
	NewAssisterFormRouter,
	domain.NewAssisterFormService,
)
