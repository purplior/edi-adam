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
	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/application/middleware"
	"github.com/podossaem/podoroot/application/router"
	"github.com/podossaem/podoroot/domain"
	"github.com/podossaem/podoroot/infra"
	"github.com/podossaem/podoroot/infra/database"
)

func StartApplication(
	databaseManager database.DatabaseManager,
	router router.Router,
) error {
	// 데이터베이스 연결
	if err := databaseManager.Init(); err != nil {
		log.Println("[#] 데이터베이스를 초기화 하는데 실패 했어요")
		return err
	}

	// 서버 생성
	app := echo.New()
	app.Use(middleware.New()...)
	router.Attach(app)

	// shutdown 채널 생성
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Start(fmt.Sprintf(":%d", config.AppPort())); err != nil {
			log.Println("[#] 서버를 시작 하는데 실패 했어요")
			panic(err)
		}
	}()

	go func() {
		databaseManager.Monitor()
	}()

	sig := <-sigChan
	log.Println("[#] 종료 시그널을 받았어요: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Println("[#] 서버를 종료 하는데 실패 했어요")
		return err
	}

	if err := databaseManager.Dispose(); err != nil {
		log.Println("[#] 데이터베이스를 종료 하는데 실패 했어요")
		return err
	}

	return nil
}

func Start() error {
	panic(
		wire.Build(
			StartApplication,
			router.New,
			domain.New,
			infra.New,
		),
	)
}
