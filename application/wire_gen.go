// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package application

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	controller2 "github.com/purplior/edi-adam/application/admin/controller"
	router2 "github.com/purplior/edi-adam/application/admin/router"
	"github.com/purplior/edi-adam/application/api/controller"
	"github.com/purplior/edi-adam/application/api/router"
	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/application/middleware"
	"github.com/purplior/edi-adam/domain/assistant"
	"github.com/purplior/edi-adam/domain/assister"
	"github.com/purplior/edi-adam/domain/auth"
	"github.com/purplior/edi-adam/domain/bookmark"
	"github.com/purplior/edi-adam/domain/category"
	"github.com/purplior/edi-adam/domain/customervoice"
	"github.com/purplior/edi-adam/domain/mission"
	"github.com/purplior/edi-adam/domain/missionlog"
	"github.com/purplior/edi-adam/domain/review"
	"github.com/purplior/edi-adam/domain/user"
	"github.com/purplior/edi-adam/domain/verification"
	"github.com/purplior/edi-adam/domain/wallet"
	"github.com/purplior/edi-adam/domain/walletlog"
	"github.com/purplior/edi-adam/infra/database"
	"github.com/purplior/edi-adam/infra/database/dynamo"
	"github.com/purplior/edi-adam/infra/database/postgre"
	"github.com/purplior/edi-adam/infra/port/openai"
	"github.com/purplior/edi-adam/infra/port/slack"
	"github.com/purplior/edi-adam/infra/port/sms"
	"github.com/purplior/edi-adam/infra/repository"
	"github.com/purplior/edi-adam/infra/session"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Injectors from app.go:

func Start() error {
	client := postgre.NewClient()
	sessionFactory := session.NewFactory(client)
	assistantRepository := repository.NewAssistantRepository(client)
	assistantService := assistant.NewAssistantService(assistantRepository)
	openaiClient := openai.NewClient()
	walletRepository := repository.NewWalletRepository(client)
	dynamoClient := dynamo.NewClient()
	walletLogRepository := repository.NewWalletLogRepository(dynamoClient)
	walletLogService := walletlog.NewWalletLogService(walletLogRepository)
	walletService := wallet.NewWalletService(walletRepository, walletLogService)
	assisterRepository := repository.NewAssisterRepository(dynamoClient)
	assisterService := assister.NewAssisterService(openaiClient, walletService, assisterRepository)
	assistantController := controller.NewAssistantController(assistantService, assisterService)
	assisterController := controller.NewAssisterController(assisterService)
	smsClient := sms.NewClient()
	verificationRepository := repository.NewVerificationRepository(client)
	verificationService := verification.NewVerificationService(smsClient, verificationRepository)
	userRepository := repository.NewUserRepository(client)
	userService := user.NewUserService(userRepository)
	missionLogRepository := repository.NewMissionLogRepository(client)
	missionLogService := missionlog.NewMissionLogService(missionLogRepository)
	authService := auth.NewAuthService(verificationService, userService, walletService, missionLogService)
	authController := controller.NewAuthController(authService)
	bookmarkRepository := repository.NewBookmarkRepository(client)
	bookmarkService := bookmark.NewBookmarkService(bookmarkRepository)
	bookmarkController := controller.NewBookmarkController(bookmarkService)
	categoryRepository := repository.NewCategoryRepository(client)
	categoryService := category.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)
	slackClient := slack.NewClient()
	customerVoiceRepository := repository.NewCustomerVoiceRepository(client)
	customerVoiceService := customervoice.NewCustomerVoiceService(slackClient, customerVoiceRepository, userService)
	customerVoiceController := controller.NewCustomerVoiceController(customerVoiceService)
	missionRepository := repository.NewMissionRepository(client)
	missionService := mission.NewMissionService(missionRepository)
	missionController := controller.NewMissionController(missionService)
	missionLogController := controller.NewMissionLogController(missionLogService)
	reviewRepository := repository.NewReviewRepository(client)
	reviewService := review.NewReviewService(reviewRepository)
	reviewController := controller.NewReviewController(reviewService)
	userController := controller.NewUserController(userService, walletService)
	verificationController := controller.NewVerificationController(verificationService)
	routerRouter := router.NewRouter(sessionFactory, assistantController, assisterController, authController, bookmarkController, categoryController, customerVoiceController, missionController, missionLogController, reviewController, userController, verificationController)
	assistantAdminController := controller2.NewAssistantAdminController(assistantService)
	router3 := router2.NewRouter(sessionFactory, assistantAdminController)
	databaseManager := database.NewDatabaseManager(client, dynamoClient)
	error2 := StartApplication(routerRouter, router3, databaseManager)
	return error2
}

// app.go:

func StartApplication(
	apiRouter router.Router,
	adminRouter router2.Router,
	dbManager database.DatabaseManager,
) error {

	if err := dbManager.Init(); err != nil {
		log.Println("[#] 데이터베이스를 초기화 하는데 실패 했어요")
		return err
	}

	app := echo.New()
	app.Use(middleware.New()...)
	app.GET("/healthx", func(c echo.Context) error {
		return c.String(200, "Hi, i'm edi-adam.")
	})

	adminGroup := app.Group("/__admin__")
	adminRouter.Attach(adminGroup)

	apiGroup := app.Group("/api/:version")
	apiRouter.Attach(apiGroup)

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
