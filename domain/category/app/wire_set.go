package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/sbec/domain/category"
)

var New = wire.NewSet(
	NewCategoryController,
	NewCategoryRouter,
	domain.NewCategoryService,
)
