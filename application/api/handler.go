package api

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/user"
)

type (
	HandlerFunc func(ctx *Context) error

	HandlerFuncOption struct {
		AdminOnly bool
	}
)

func Handler(handler HandlerFunc, opt HandlerFuncOption) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		appCtx, ok := ctx.(*Context)
		if !ok {
			appCtx = &Context{
				Context: ctx,
			}
		}
		if opt.AdminOnly {
			isNotAllowed := appCtx.Identity == nil || appCtx.Identity.Role != user.Role_Master
			if isNotAllowed {
				return appCtx.SendError(exception.ErrUnauthorized)
			}
		}

		return handler(appCtx)
	}
}
