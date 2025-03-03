package controller

import "github.com/google/wire"

var New = wire.NewSet(
	NewAssistantController,
	NewAssisterController,
	NewAuthController,
	NewBookmarkController,
	NewCategoryController,
	NewCustomerVoiceController,
	NewMissionController,
	NewMissionLogController,
	NewReviewController,
	NewUserController,
	NewVerificationController,
)
