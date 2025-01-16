package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/review"
)

var New = wire.NewSet(
	NewReviewRouter,
	NewReviewController,
	domain.NewReviewService,
)
