package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/challenge"
)

var New = wire.NewSet(
	NewChallengeController,
	NewChallengeRouter,
	domain.NewChallengeService,
)
