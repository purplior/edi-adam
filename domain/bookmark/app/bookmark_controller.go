package app

import (
	domain "github.com/purplior/sbec/domain/bookmark"
	"github.com/purplior/sbec/domain/shared/inner"
)

type (
	BookmarkController interface {
	}
)

type (
	bookmarkController struct {
		bookmarkService domain.BookmarkService
		cm              inner.ContextManager
	}
)

func NewBookmarkController(
	bookmarkService domain.BookmarkService,
	cm inner.ContextManager,
) BookmarkController {
	return &bookmarkController{
		bookmarkService: bookmarkService,
		cm:              cm,
	}
}
