package app

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/domain/verification"
)

var New = wire.NewSet(
	NewEmailVerificationController,
	NewVerificationRouter,
	verification.NewEmailVerificationService,
)
