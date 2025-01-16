package repository

import "github.com/google/wire"

var New = wire.NewSet(
	NewAssistantRepository,
	NewAssisterRepository,
	NewBookmarkRepository,
	NewCategoryRepository,
	NewChallengeRepository,
	NewCustomerVoiceRepository,
	NewEmailVerificationRepository,
	NewLedgerRepository,
	NewMissionRepository,
	NewPhoneVerificationRepository,
	NewReviewRepository,
	NewUserRepository,
	NewWalletRepository,
)
