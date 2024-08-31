package app

import (
	"github.com/google/wire"
	"github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/domain/verification/persist"
)

var New = wire.NewSet(
	NewEmailVerificationController,
	NewRouter,
	persist.NewEmailVerificationRepository,
	verification.NewEmailVerificationService,
)
