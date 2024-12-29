package app

import (
	"github.com/google/wire"
	"github.com/purplior/podoroot/domain/verification"
)

var New = wire.NewSet(
	NewEmailVerificationController,
	NewPhoneVerificationController,
	NewVerificationRouter,
	verification.NewEmailVerificationService,
	verification.NewPhoneVerificationService,
)
