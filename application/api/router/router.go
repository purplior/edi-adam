package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/application/api/controller"
	"github.com/purplior/edi-adam/application/common"
	"github.com/purplior/edi-adam/domain/shared/inner"
)

type (
	Router interface {
		Attach(router *echo.Group)
	}
)

type (
	apiRouter struct {
		controllers []controller.Controller
		hf          common.HandlerFactory
	}
)

func (r *apiRouter) Attach(router *echo.Group) {
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
	assistantController controller.AssistantController,
	assisterController controller.AssisterController,
	authController controller.AuthController,
	bookmarkController controller.BookmarkController,
	categoryController controller.CategoryController,
	customerVoiceController controller.CustomerVoiceController,
	missionController controller.MissionController,
	missionLogController controller.MissionLogController,
	reviewController controller.ReviewController,
	userController controller.UserController,
	verificationController controller.VerificationController,
) Router {
	return &apiRouter{
		hf: common.NewHandlerFactory(sessionFactory),
		controllers: []controller.Controller{
			assistantController,
			assisterController,
			authController,
			bookmarkController,
			categoryController,
			customerVoiceController,
			missionController,
			missionLogController,
			reviewController,
			userController,
			verificationController,
		},
	}
}
