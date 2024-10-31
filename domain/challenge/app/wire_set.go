package app

import (
	"github.com/google/wire"
	domain "github.com/podossaem/podoroot/domain/challenge"
)

var New = wire.NewSet(
	NewChallengeController,
	NewChallengeRouter,
	domain.NewChallengeService,
)
