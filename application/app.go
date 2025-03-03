//go:build wireinject
// +build wireinject

package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/purplior/edi-adam/application/admin"
	adminrouter "github.com/purplior/edi-adam/application/admin/router"
	"github.com/purplior/edi-adam/application/api"
	apirouter "github.com/purplior/edi-adam/application/api/router"
	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/application/middleware"
	"github.com/purplior/edi-adam/domain"
	"github.com/purplior/edi-adam/infra"
	"github.com/purplior/edi-adam/infra/database"
)

func StartApplication(
	apiRouter apirouter.Router,
	adminRouter adminrouter.Router,
	dbManager database.DatabaseManager,
) error {
	// 데이터베이스 연결
	if err := dbManager.Init(); err != nil {
		log.Println("[#] 데이터베이스를 초기화 하는데 실패 했어요")
		return err
	}

	// 서버 생성
	app := echo.New()
	app.Use(middleware.New()...)
	app.GET("/healthx", func(c echo.Context) error {
		return c.String(200, "Hi, i'm edi-adam.")
	})

	adminGroup := app.Group("/__admin__")
	adminRouter.Attach(adminGroup)

	apiGroup := app.Group("/api/:version")
	apiRouter.Attach(apiGroup)

	// shutdown 채널 생성
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Start(fmt.Sprintf(":%d", config.Port())); err != nil {
			log.Println("[#] 서버를 시작 하는데 실패 했어요")
			panic(err)
		}
	}()

	go func() {
		dbManager.Monitor()
	}()

	sig := <-sigChan
	log.Println("[#] 종료 시그널을 받았어요: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Println("[#] 서버를 종료 하는데 실패 했어요")
		return err
	}

	if err := dbManager.Dispose(); err != nil {
		log.Println("[#] 데이터베이스를 종료 하는데 실패 했어요")
		return err
	}

	return nil
}

func Start() error {
	panic(
		wire.Build(
			StartApplication,
			infra.New,
			admin.New,
			api.New,
			domain.New,
		),
	)
}
