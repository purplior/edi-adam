package repository

import "github.com/google/wire"

var New = wire.NewSet(
	NewAssistantRepository,
	NewAssisterFormRepository,
	NewEmailVerificationRepository,
	NewUserRepository,
)
