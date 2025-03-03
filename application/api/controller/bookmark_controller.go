package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/bookmark"
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	BookmarkController interface {
		Controller
		// 내 북마크 보기
		GetMine() common.Route
		// 내 북마크 존재하는지 보기
		GetMine_PaginatedList() common.Route
		// 내 북마크 토글하기
		ToggleMine() common.Route
	}
)

type (
	bookmarkController struct {
		bookmarkService domain.BookmarkService
	}
)

func (c *bookmarkController) GroupPath() string {
	return "/bookmarks"
}

func (c *bookmarkController) Routes() []common.Route {
	return []common.Route{
		c.GetMine(),
		c.GetMine_PaginatedList(),
		c.ToggleMine(),
	}
}

func (c *bookmarkController) GetMine() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/m/fst",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			AssistantID uint `query:"assistant_id"`
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		m, err := c.bookmarkService.Get(
			ctx.Session(),
			domain.QueryOption{
				UserID:      ctx.Identity.ID,
				AssistantID: dto.AssistantID,
			},
		)
		if err != nil && err != exception.ErrNoRecord {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Bookmark model.Bookmark `json:"bookmark"`
		}{
			Bookmark: m,
		})
	}

	return route
}

func (c *bookmarkController) GetMine_PaginatedList() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/m/lst",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			pagination.PaginationRequest
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		mArr, pageMeta, err := c.bookmarkService.GetPaginatedList(
			ctx.Session(),
			pagination.PaginationQuery[domain.QueryOption]{
				QueryOption: domain.QueryOption{
					UserID: ctx.Identity.ID,
				},
				PageRequest: dto.PaginationRequest,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Bookmarks []model.Bookmark          `json:"bookmarks"`
			PageMeta  pagination.PaginationMeta `json:"pageMeta"`
		}{
			Bookmarks: mArr,
			PageMeta:  pageMeta,
		})
	}

	return route
}

func (c *bookmarkController) ToggleMine() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Path:   "/m/fst/toggle",
		Option: common.RouteOption{Member: true},
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			AssistantID uint `body:"assistantId,required"`
			Toggle      bool `body:"toggle"`
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		if err = c.bookmarkService.Toggle(
			ctx.Session(),
			domain.QueryOption{
				UserID:      ctx.Identity.ID,
				AssistantID: dto.AssistantID,
			},
			dto.Toggle,
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(nil)
	}

	return route
}

func NewBookmarkController(
	bookmarkService domain.BookmarkService,
) BookmarkController {
	return &bookmarkController{
		bookmarkService: bookmarkService,
	}
}
