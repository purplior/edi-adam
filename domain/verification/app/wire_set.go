package app

import (
	"github.com/google/wire"
	"github.com/podossaem/root/domain/verification"
	"github.com/podossaem/root/domain/verification/persist"
)

var New = wire.NewSet(
	NewEmailVerificationController,
	NewRouter,
	persist.NewEmailVerificationRepository,
	verification.NewEmailVerificationService,
)
