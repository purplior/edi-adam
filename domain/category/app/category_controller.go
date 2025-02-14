package app

import (
	"github.com/purplior/sbec/application/api"
	"github.com/purplior/sbec/application/response"
	domain "github.com/purplior/sbec/domain/category"
	"github.com/purplior/sbec/domain/shared/inner"
)

type (
	CategoryController interface {
		GetMainCategoryInfoList() api.HandlerFunc
	}
)

type (
	categoryController struct {
		categoryService domain.CategoryService
		cm              inner.ContextManager
	}
)

func (c *categoryController) GetMainCategoryInfoList() api.HandlerFunc {
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
