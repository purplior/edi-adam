package app

import (
	domain "github.com/purplior/podoroot/domain/bookmark"
	"github.com/purplior/podoroot/domain/shared/inner"
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
