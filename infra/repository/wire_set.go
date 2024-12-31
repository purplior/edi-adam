package repository

import "github.com/google/wire"

var New = wire.NewSet(
	NewAssistantRepository,
	NewAssisterRepository,
	NewAssisterFormRepository,
	NewBookmarkRepository,
	NewCategoryRepository,
	NewChallengeRepository,
	NewCustomerVoiceRepository,
	NewEmailVerificationRepository,
	NewLedgerRepository,
	NewMissionRepository,
	NewPhoneVerificationRepository,
	NewUserRepository,
	NewWalletRepository,
)
