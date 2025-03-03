package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/category"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	CategoryController interface {
		Controller
		// 주요 카테고리 목록 가져오기
		GetMainList() common.Route
	}
)

type (
	categoryController struct {
		categoryService domain.CategoryService
	}
)

func (c *categoryController) GroupPath() string {
	return "/categories"
}

func (c *categoryController) Routes() []common.Route {
	return []common.Route{
		c.GetMainList(),
	}
}

func (c *categoryController) GetMainList() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/o/lst/main",
	}
	route.Handler = func(ctx *common.Context) error {
		categories, err := c.categoryService.GetMainList(ctx.Session())
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Categories []model.Category `json:"categories"`
		}{
			Categories: categories,
		})
	}

	return route
}

func NewCategoryController(
	categoryService domain.CategoryService,
) CategoryController {
	return &categoryController{
		categoryService: categoryService,
	}
}
