package api

import (
	"github.com/labstack/echo/v4"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

type (
	Router interface {
		Attach(app *echo.Echo)
	}

	router struct {
		verificationRouter verification.Router
	}
)

func (r *router) Attach(app *echo.Echo) {
	api := app.Group("/api/:version")

	r.verificationRouter.Attach(api)
}

func NewRouter(
	verificationRouter verification.Router,
) Router {
	return &router{
		verificationRouter: verificationRouter,
	}
}
