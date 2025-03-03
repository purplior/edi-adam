package common

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	HandlerFunc func(ctx *Context) error

	RouteOption struct {
		Member  bool
		Admin   bool
		Timeout time.Duration
	}

	Route struct {
		Method  string
		Path    string
		Option  RouteOption
		Handler HandlerFunc
	}
)

type (
	HandlerFactory = func(HandlerFunc, RouteOption) echo.HandlerFunc
)

// // @컨트롤러 핸들러 레벨의 옵셔널 처리를 하는 미들웨어 생성자
func NewHandlerFactory(
	sessionFactory inner.SessionFactory,
) HandlerFactory {
	return func(hf HandlerFunc, ro RouteOption) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, ok := c.(*Context)
			if !ok {
				ctx = &Context{
					Context: c,
				}
			}

			if ro.Admin {
				isNotAllowed := ctx.Identity == nil || ctx.Identity.Role != model.UserRole_Admin
				if isNotAllowed {
					return ctx.SendError(exception.ErrUnauthorized)
				}
			} else if ro.Member {
				isNotAllowed := ctx.Identity == nil || ctx.Identity.Role != model.UserRole_Member
				if isNotAllowed {
					return ctx.SendError(exception.ErrUnauthorized)
				}
			}

			timeout := ro.Timeout
			if timeout <= 0 {
				timeout = time.Duration(12 * time.Second)
			}
			session, cancel := sessionFactory.CreateNewSession(timeout)
			defer cancel()
			session.SetIdentity(ctx.Identity)
			ctx.session = session

			return hf(ctx)
		}
	}
}
