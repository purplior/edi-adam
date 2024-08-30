//go:build wireinject
// +build wireinject

package application

import (
	"fmt"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podossaem/root/application/api"
	"github.com/podossaem/root/application/config"
	"github.com/podossaem/root/domain"
)

func StartApplication(
	router api.Router,
) error {
	app := echo.New()

	app.Use(middleware.Logger())
	router.Attach(app)

	return app.Start(fmt.Sprintf(":%d", config.AppPort()))
}

func Start() error {
	panic(
		wire.Build(
			StartApplication,
			api.NewRouter,
			domain.New,
		),
	)
}
