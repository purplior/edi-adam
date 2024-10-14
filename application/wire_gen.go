// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package application

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/application/middleware"
	"github.com/podossaem/podoroot/application/router"
	"github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/assistant/app"
	"github.com/podossaem/podoroot/domain/assistant/persist"
	"github.com/podossaem/podoroot/domain/auth"
	app2 "github.com/podossaem/podoroot/domain/auth/app"
	"github.com/podossaem/podoroot/domain/me"
	app3 "github.com/podossaem/podoroot/domain/me/app"
	"github.com/podossaem/podoroot/domain/user"
	app4 "github.com/podossaem/podoroot/domain/user/app"
	persist3 "github.com/podossaem/podoroot/domain/user/persist"
	"github.com/podossaem/podoroot/domain/verification"
	app5 "github.com/podossaem/podoroot/domain/verification/app"
	persist2 "github.com/podossaem/podoroot/domain/verification/persist"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Injectors from app.go:

func Start() error {
	client := podosql.NewClient()
	databaseManager := database.NewDatabaseManager(client)
	assistantRepository := persist.NewAssistantRepository(client)
	assistantService := assistant.NewAssistantService(assistantRepository)
	assistantController := app.NewAssistantController(assistantService)
	assistantRouter := app.NewAssistantRouter(assistantController)
	emailVerificationRepository := persist2.NewEmailVerificationRepository(client)
	emailVerificationService := verification.NewEmailVerificationService(emailVerificationRepository)
	userRepository := persist3.NewUserRepository(client)
	userService := user.NewUserService(userRepository)
	authService := auth.NewAuthService(emailVerificationService, userService)
	authController := app2.NewAuthController(authService)
	authRouter := app2.NewAuthRouter(authController)
	meService := me.NewMeService(userRepository)
	meController := app3.NewMeController(meService)
	meRouter := app3.NewMeRouter(meController)
	userController := app4.NewUserController()
	userRouter := app4.NewUserRouter(userController)
	emailVerificationController := app5.NewEmailVerificationController(emailVerificationService)
	verificationRouter := app5.NewVerificationRouter(emailVerificationController)
	routerRouter := router.New(assistantRouter, authRouter, meRouter, userRouter, verificationRouter)
	error2 := StartApplication(databaseManager, routerRouter)
	return error2
}

// app.go:

func StartApplication(
	databaseManager database.DatabaseManager, router2 router.Router,

) error {

	if err := databaseManager.Init(); err != nil {
		log.Println("[#] 데이터베이스를 초기화 하는데 실패 했어요")
		return err
	}
	app6 := echo.New()
	app6.
		Use(middleware.New()...)
	router2.
		Attach(app6)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app6.Start(fmt.Sprintf(":%d", config.AppPort())); err != nil {
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

	if err := app6.Shutdown(ctx); err != nil {
		log.Println("[#] 서버를 종료 하는데 실패 했어요")
		return err
	}

	if err := databaseManager.Dispose(); err != nil {
		log.Println("[#] 데이터베이스를 종료 하는데 실패 했어요")
		return err
	}

	return nil
}
