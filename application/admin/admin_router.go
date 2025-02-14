package admin

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/sbec/application/api"
)

type (
	AdminRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	adminRouter struct {
		adminController AdminController
	}
)

func (r *adminRouter) Attach(router *echo.Group) {
	assistantRouter := router.Group("/assistants")

	assistantRouter.POST(
		"/approve-one",
		api.Handler(
			r.adminController.ApproveAssistantOne(),
			api.HandlerFuncOption{AdminOnly: true},
		),
	)
}

func NewAdminRouter(
	adminController AdminController,
) AdminRouter {
	return &adminRouter{
		adminController: adminController,
	}
}
