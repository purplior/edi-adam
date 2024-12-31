// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package application

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/application/middleware"
	"github.com/purplior/podoroot/application/router"
	"github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/assistant/app"
	"github.com/purplior/podoroot/domain/assister"
	app2 "github.com/purplior/podoroot/domain/assister/app"
	"github.com/purplior/podoroot/domain/assisterform"
	app3 "github.com/purplior/podoroot/domain/assisterform/app"
	"github.com/purplior/podoroot/domain/auth"
	app4 "github.com/purplior/podoroot/domain/auth/app"
	"github.com/purplior/podoroot/domain/bookmark"
	app5 "github.com/purplior/podoroot/domain/bookmark/app"
	"github.com/purplior/podoroot/domain/category"
	app6 "github.com/purplior/podoroot/domain/category/app"
	"github.com/purplior/podoroot/domain/challenge"
	app7 "github.com/purplior/podoroot/domain/challenge/app"
	"github.com/purplior/podoroot/domain/customervoice"
	app8 "github.com/purplior/podoroot/domain/customervoice/app"
	"github.com/purplior/podoroot/domain/ledger"
	"github.com/purplior/podoroot/domain/me"
	app9 "github.com/purplior/podoroot/domain/me/app"
	"github.com/purplior/podoroot/domain/mission"
	app10 "github.com/purplior/podoroot/domain/mission/app"
	"github.com/purplior/podoroot/domain/user"
	app11 "github.com/purplior/podoroot/domain/user/app"
	"github.com/purplior/podoroot/domain/verification"
	app12 "github.com/purplior/podoroot/domain/verification/app"
	"github.com/purplior/podoroot/domain/wallet"
	"github.com/purplior/podoroot/infra"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podomongo"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/port/podoopenai"
	"github.com/purplior/podoroot/infra/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Injectors from app.go:

func Start() error {
	client := podosql.NewClient()
	podomongoClient := podomongo.NewClient()
	databaseManager := database.NewDatabaseManager(client, podomongoClient)
	assistantRepository := repository.NewAssistantRepository(client)
	podoopenaiClient := podoopenai.NewClient()
	assisterFormRepository := repository.NewAssisterFormRepository(podomongoClient)
	assisterFormService := assisterform.NewAssisterFormService(assisterFormRepository)
	walletRepository := repository.NewWalletRepository(client)
	ledgerRepository := repository.NewLedgerRepository(client)
	ledgerService := ledger.NewLedgerService(ledgerRepository)
	walletService := wallet.NewWalletService(walletRepository, ledgerService)
	assisterRepository := repository.NewAssisterRepository(client)
	contextManager := infra.NewContextManager(client)
	assisterService := assister.NewAssisterService(podoopenaiClient, assisterFormService, walletService, assisterRepository, contextManager)
	assistantService := assistant.NewAssistantService(assistantRepository, assisterService, assisterFormService, contextManager)
	assistantController := app.NewAssistantController(assistantService, contextManager)
	assistantRouter := app.NewAssistantRouter(assistantController)
	assisterController := app2.NewAssisterController(assisterService, contextManager)
	assisterRouter := app2.NewAssisterRouter(assisterController)
	assisterFormController := app3.NewAssisterFormController(assisterFormService, contextManager)
	assisterFormRouter := app3.NewAssisterFormRouter(assisterFormController)
	emailVerificationRepository := repository.NewEmailVerificationRepository(client)
	emailVerificationService := verification.NewEmailVerificationService(emailVerificationRepository)
	phoneVerificationRepository := repository.NewPhoneVerificationRepository(client)
	phoneVerificationService := verification.NewPhoneVerificationService(phoneVerificationRepository)
	userRepository := repository.NewUserRepository(client)
	userService := user.NewUserService(userRepository)
	challengeRepository := repository.NewChallengeRepository(client)
	challengeService := challenge.NewChallengeService(challengeRepository, contextManager)
	authService := auth.NewAuthService(emailVerificationService, phoneVerificationService, userService, walletService, challengeService, contextManager)
	authController := app4.NewAuthController(authService, contextManager)
	authRouter := app4.NewAuthRouter(authController)
	bookmarkRepository := repository.NewBookmarkRepository(client)
	bookmarkService := bookmark.NewBookmarkService(bookmarkRepository)
	bookmarkController := app5.NewBookmarkController(bookmarkService, contextManager)
	bookmarkRouter := app5.NewBookmarkRouter(bookmarkController)
	categoryRepository := repository.NewCategoryRepository(client)
	categoryService := category.NewCategoryService(categoryRepository)
	categoryController := app6.NewCategoryController(categoryService, contextManager)
	categoryRouter := app6.NewCategoryRouter(categoryController)
	challengeController := app7.NewChallengeController(challengeService, contextManager)
	challengeRouter := app7.NewChallengeRouter(challengeController)
	customerVoiceRepository := repository.NewCustomerVoiceRepository(client)
	customerVoiceService := customervoice.NewCustomerVoiceService(customerVoiceRepository, userService, contextManager)
	customerVoiceController := app8.NewCustomerVoiceController(customerVoiceService, contextManager)
	customerVoiceRouter := app8.NewCustomerVoiceRouter(customerVoiceController)
	meService := me.NewMeService()
	meController := app9.NewMeController(meService, assistantService, assisterFormService, authService, userService, walletService, bookmarkService, contextManager)
	meRouter := app9.NewMeRouter(meController)
	missionRepository := repository.NewMissionRepository(client)
	missionService := mission.NewMissionService(missionRepository, challengeService, walletService, contextManager)
	missionController := app10.NewMissionController(missionService, contextManager)
	missionRouter := app10.NewMissionRouter(missionController)
	userController := app11.NewUserController(userService, contextManager)
	userRouter := app11.NewUserRouter(userController)
	emailVerificationController := app12.NewEmailVerificationController(emailVerificationService, contextManager)
	phoneVerificationController := app12.NewPhoneVerificationController(phoneVerificationService, userService, contextManager)
	verificationRouter := app12.NewVerificationRouter(emailVerificationController, phoneVerificationController)
	routerRouter := router.New(assistantRouter, assisterRouter, assisterFormRouter, authRouter, bookmarkRouter, categoryRouter, challengeRouter, customerVoiceRouter, meRouter, missionRouter, userRouter, verificationRouter)
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
	app13 := echo.New()
	app13.
		Use(middleware.New()...)
	router2.
		Attach(app13)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app13.Start(fmt.Sprintf(":%d", config.AppPort())); err != nil {
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

	if err := app13.Shutdown(ctx); err != nil {
		log.Println("[#] 서버를 종료 하는데 실패 했어요")
		return err
	}

	if err := databaseManager.Dispose(); err != nil {
		log.Println("[#] 데이터베이스를 종료 하는데 실패 했어요")
		return err
	}

	return nil
}
