package app

import "github.com/labstack/echo/v4"

type (
	BookmarkRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	bookmarkRouter struct {
		bookmarkController BookmarkController
	}
)

func (r *bookmarkRouter) Attach(router *echo.Group) {
}

func NewBookmarkRouter(
	bookmarkController BookmarkController,
) BookmarkRouter {
	return &bookmarkRouter{
		bookmarkController: bookmarkController,
	}
}
