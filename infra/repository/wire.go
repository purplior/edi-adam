package repository

import (
	"github.com/google/wire"
)

var New = wire.NewSet(
	NewAssistantRepository,
	NewAssisterRepository,
	NewBookmarkRepository,
	NewCategoryRepository,
	NewCustomerVoiceRepository,
	NewMissionRepository,
	NewMissionLogRepository,
	NewReviewRepository,
	NewUserRepository,
	NewVerificationRepository,
	NewWalletRepository,
	NewWalletLogRepository,
)
