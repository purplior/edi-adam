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
	"github.com/podossaem/root/infra"
	"github.com/podossaem/root/infra/database"
	"github.com/podossaem/root/infra/database/mymongo"
)

func StartApplication(
	router api.Router,
	mymongoClient *mymongo.Client,
) error {
	if err := database.Init(mymongoClient); err != nil {
		return err
	}

	app := echo.New()

	app.Use(middleware.Logger())
	router.Attach(app)

	if err := app.Start(fmt.Sprintf(":%d", config.AppPort())); err != nil {
		return err
	}

	defer func() {
		if err := database.Dispose(mymongoClient); err != nil {
			panic(err)
		}
	}()

	return nil
}

func Start() error {
	panic(
		wire.Build(
			StartApplication,
			api.NewRouter,
			domain.New,
			infra.New,
		),
	)
}
