package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/sbec/application/api"
)

type (
	CustomerVoiceRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	customerVoiceRouter struct {
		customerVoiceController CustomerVoiceController
	}
)

func (r *customerVoiceRouter) Attach(router *echo.Group) {
	customerVoiceRouterGroup := router.Group("/customer_voices")

	customerVoiceRouterGroup.POST(
		"/",
		api.Handler(
			r.customerVoiceController.RegisterOne(),
			api.HandlerFuncOption{},
		),
	)
}

func NewCustomerVoiceRouter(
	customerVoiceController CustomerVoiceController,
) CustomerVoiceRouter {
	return &customerVoiceRouter{
		customerVoiceController: customerVoiceController,
	}
}
