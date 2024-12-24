package app

import (
	"github.com/purplior/podoroot/application/api"
	"github.com/purplior/podoroot/application/response"
	domain "github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/domain/shared/inner"
)

type (
	CategoryController interface {
		GetMainCategoryInfos() api.HandlerFunc
	}
)

type (
	categoryController struct {
		categoryService domain.CategoryService
		cm              inner.ContextManager
	}
)

func (c *categoryController) GetMainCategoryInfos() api.HandlerFunc {
	return func(ctx *api.Context) error {
		innerCtx, cancel := c.cm.NewContext()
		defer cancel()

		categoryList, err := c.categoryService.GetMainCategoryList(innerCtx)
		if err != nil {
			return ctx.SendError(err)
		}

		categoryInfos := make([]domain.CategoryInfo, len(categoryList))
		for i, category := range categoryList {
			categoryInfos[i] = category.ToInfo()
		}

		return ctx.SendJSON(response.JSONResponse{
			Data: struct {
				CategoryInfos []domain.CategoryInfo `json:"categoryInfos"`
			}{
				CategoryInfos: categoryInfos,
			},
		})
	}
}

func NewCategoryController(
	categoryService domain.CategoryService,
	cm inner.ContextManager,
) CategoryController {
	return &categoryController{
		categoryService: categoryService,
		cm:              cm,
	}
}
