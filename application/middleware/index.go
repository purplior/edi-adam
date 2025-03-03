package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Logger(),

		// TODO: alpha 서버에서만 *로 설정 필요
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}),

		NewContextMiddleware(),
		NewAuthMiddleware(),
	}
}
