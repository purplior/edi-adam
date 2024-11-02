package repository

import "github.com/google/wire"

var New = wire.NewSet(
	NewAssistantRepository,
	NewAssisterRepository,
	NewAssisterFormRepository,
	NewChallengeRepository,
	NewEmailVerificationRepository,
	NewLedgerRepository,
	NewMissionRepository,
	NewUserRepository,
	NewWalletRepository,
)
