package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/root/application/api/controller"
)

type (
	Router interface {
		Attach(router *echo.Group)
	}

	router struct {
		controller EmailVerificationController
	}
)

func (r *router) Attach(router *echo.Group) {
	emailVerificationRouter := router.Group("/email-verifications")
	emailVerificationRouter.POST("/request-code", controller.Handler(r.controller.RequestCode()))
}

func NewRouter(
	controller EmailVerificationController,
) Router {
	return &router{
		controller: controller,
	}
}
