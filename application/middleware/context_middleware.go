package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/logger"
)

// @MIDDLEWARE {컨텍스트 미들웨어}
// - API 컨텍스트 생성
// - API 트랜잭션 내에서 panic이 일어날 때 JSON 응답 처리
func NewContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiContext := &common.Context{
				Context: c,
			}

			defer apiContext.ClearSession()
			defer func() {
				if r := recover(); r != nil {
					if !apiContext.Response().Committed {
						logger.DebugAny(r)
						err := exception.ErrInternalServer
						if rErr, ok := r.(error); ok {
							err = rErr
						}
						apiContext.SendError(err)
					}
				}
			}()

			return next(apiContext)
		}
	}
}
