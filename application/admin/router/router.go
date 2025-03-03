package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/application/admin/controller"
	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/domain/shared/inner"
)

type (
	Router interface {
		Attach(router *echo.Group)
	}
)

type (
	adminRouter struct {
		controllers []controller.Controller
		hf          common.HandlerFactory
	}
)

func (r *adminRouter) Attach(router *echo.Group) {
	for _, controller := range r.controllers {
		rg := router.Group(controller.GroupPath())
		routes := controller.Routes()

		for _, route := range routes {
			switch route.Method {
			case http.MethodGet:
				rg.GET(route.Path, r.hf(route.Handler, route.Option))
			case http.MethodPost:
				rg.POST(route.Path, r.hf(route.Handler, route.Option))
			case http.MethodPatch:
				rg.PATCH(route.Path, r.hf(route.Handler, route.Option))
			case http.MethodPut:
				rg.PUT(route.Path, r.hf(route.Handler, route.Option))
			case http.MethodDelete:
				rg.DELETE(route.Path, r.hf(route.Handler, route.Option))
			}
		}
	}
}

func NewRouter(
	sessionFactory inner.SessionFactory,
	assistantAdminController controller.AssistantAdminController,
) Router {
	return &adminRouter{
		hf: common.NewHandlerFactory(sessionFactory),
		controllers: []controller.Controller{
			assistantAdminController,
		},
	}
}
