package app

import (
	"github.com/google/wire"
	domain "github.com/purplior/podoroot/domain/category"
)

var New = wire.NewSet(
	NewCategoryController,
	NewCategoryRouter,
	domain.NewCategoryService,
)
