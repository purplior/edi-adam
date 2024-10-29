package repository

import "github.com/google/wire"

var New = wire.NewSet(
	NewAssistantRepository,
	NewAssisterRepository,
	NewAssisterFormRepository,
	NewEmailVerificationRepository,
	NewLedgerRepository,
	NewUserRepository,
	NewWalletRepository,
)
