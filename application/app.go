//go:build wireinject
// +build wireinject

package application

import (
	"fmt"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podossaem/podoroot/application/api"
	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain"
	"github.com/podossaem/podoroot/infra"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/mymongo"
	"github.com/podossaem/podoroot/infra/database/myredis"
)

func StartApplication(
	router api.Router,
	mymongoClient *mymongo.Client,
	myredisClient *myredis.Client,
) error {
	if err := database.Init(
		mymongoClient,
		myredisClient,
	); err != nil {
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
