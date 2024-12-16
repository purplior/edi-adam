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
	"github.com/purplior/podoroot/domain/challenge"
	app5 "github.com/purplior/podoroot/domain/challenge/app"
	"github.com/purplior/podoroot/domain/ledger"
	"github.com/purplior/podoroot/domain/me"
	app6 "github.com/purplior/podoroot/domain/me/app"
	"github.com/purplior/podoroot/domain/mission"
	app7 "github.com/purplior/podoroot/domain/mission/app"
	"github.com/purplior/podoroot/domain/user"
	app8 "github.com/purplior/podoroot/domain/user/app"
	"github.com/purplior/podoroot/domain/verification"
	app9 "github.com/purplior/podoroot/domain/verification/app"
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
	userRepository := repository.NewUserRepository(client)
	userService := user.NewUserService(userRepository)
	challengeRepository := repository.NewChallengeRepository(client)
	challengeService := challenge.NewChallengeService(challengeRepository, contextManager)
	authService := auth.NewAuthService(emailVerificationService, userService, walletService, challengeService, contextManager)
	authController := app4.NewAuthController(authService, contextManager)
	authRouter := app4.NewAuthRouter(authController)
	challengeController := app5.NewChallengeController(challengeService, contextManager)
	challengeRouter := app5.NewChallengeRouter(challengeController)
	meService := me.NewMeService(userRepository)
	meController := app6.NewMeController(meService, authService, userService, walletService, contextManager)
	meRouter := app6.NewMeRouter(meController)
	missionRepository := repository.NewMissionRepository(client)
	missionService := mission.NewMissionService(missionRepository, challengeService, walletService, contextManager)
	missionController := app7.NewMissionController(missionService, contextManager)
	missionRouter := app7.NewMissionRouter(missionController)
	userController := app8.NewUserController()
	userRouter := app8.NewUserRouter(userController)
	emailVerificationController := app9.NewEmailVerificationController(emailVerificationService, contextManager)
	verificationRouter := app9.NewVerificationRouter(emailVerificationController)
	routerRouter := router.New(assistantRouter, assisterRouter, assisterFormRouter, authRouter, challengeRouter, meRouter, missionRouter, userRouter, verificationRouter)
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
	app10 := echo.New()
	app10.
		Use(middleware.New()...)
	router2.
		Attach(app10)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app10.Start(fmt.Sprintf(":%d", config.AppPort())); err != nil {
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

	if err := app10.Shutdown(ctx); err != nil {
		log.Println("[#] 서버를 종료 하는데 실패 했어요")
		return err
	}

	if err := databaseManager.Dispose(); err != nil {
		log.Println("[#] 데이터베이스를 종료 하는데 실패 했어요")
		return err
	}

	return nil
}
