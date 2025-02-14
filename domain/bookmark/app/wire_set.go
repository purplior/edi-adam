package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/bookmark"
)

var New = wire.NewSet(
	NewBookmarkRouter,
	NewBookmarkController,
	domain.NewBookmarkService,
)
