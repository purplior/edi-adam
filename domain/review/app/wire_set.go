package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/review"
)

var New = wire.NewSet(
	NewReviewRouter,
	NewReviewController,
	domain.NewReviewService,
)
