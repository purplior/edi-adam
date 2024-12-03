package app

import (
	"github.com/google/wire"
	"github.com/purplior/podoroot/domain/verification"
)

var New = wire.NewSet(
	NewEmailVerificationController,
	NewVerificationRouter,
	verification.NewEmailVerificationService,
)
