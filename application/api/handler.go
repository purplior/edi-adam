package api

import (
	"github.com/labstack/echo/v4"
)

type HandlerFunc = func(ctx *Context) error

func Handler(handler HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		appCtx, ok := ctx.(*Context)
		if !ok {
			appCtx = &Context{
				Context: ctx,
			}
		}

		return handler(appCtx)
	}
}
