package postgre

import (
	"log"
	"os"
	"time"

	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/shared/constant"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	DB = gorm.DB

	Association = gorm.Association

	ConstructorOption struct {
		Phase constant.Phase
		DSN   string
	}

	Client struct {
		*gorm.DB
		opt ConstructorOption
	}
)

func (c *Client) DBWithContext(session inner.Session) *DB {
	tx := session.Transaction()
	if tx != nil {
		return tx
	}

	// TODO: 한 세션 내에서 하나의 db context만 사용하도록 수정 필요
	return c.DB.WithContext(session.Context())
}

func (c *Client) ConnectDB() error {
	isDebugMode := config.DebugMode()

	var dbLogger logger.Interface
	if isDebugMode {
		dbLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(c.opt.DSN), &gorm.Config{
		Logger:      dbLogger,
		PrepareStmt: true,
	})
	if err != nil {
		return err
	}

	if isDebugMode {
		log.Println("[postgre] 디버그모드가 활성화 되었어요")
		db.Debug()
	}

	c.DB = db
	postgre, err := c.DB.DB()
	if err != nil {
		return err
	}

	postgre.SetMaxIdleConns(10)           // 유휴 상태로 유지할 수 있는 최대 연결 수
	postgre.SetMaxOpenConns(100)          // 데이터베이스에 열어둘 최대 연결 수
	postgre.SetConnMaxLifetime(time.Hour) // 연결의 최대 수명 (예: 1시간)

	if err := c.PingDB(); err != nil {
		return err
	}

	log.Println("[postgre] 데이터베이스 연결에 성공 했어요")

	return nil
}

func (c *Client) MigrateDB() error {
	// Relation 때문에 AutoMigrate 순서가 중요함
	return c.DB.AutoMigrate(
		model.User{},
		model.Category{},
		model.Wallet{},
		model.WalletLog{},
		model.Assistant{},
		model.Verification{},
		model.Mission{},
		model.MissionLog{},
		model.CustomerVoice{},
		model.Bookmark{},
		model.Review{},
	)
}

func (c *Client) PingDB() error {
	postgre, err := c.DB.DB()
	if err != nil {
		return err
	}

	return postgre.Ping()
}

func (c *Client) ReconnectDB(
	maxAttempts int,
	baseDelay time.Duration,
) {
	attempts := 0
	for {
		err := c.PingDB()
		if err == nil {
			log.Println("[postgre] 데이터베이스 재연결에 성공 했어요")
			return
		}

		attempts++
		if attempts > maxAttempts {

			log.Fatalf("[postgre] 데이터베이스 재시도 연결 횟수가 최대를 초과 했어요: %v", err)
			return
		}

		// 2의 지수 증가
		delay := baseDelay * time.Duration(1<<attempts)
		log.Printf("[postgre] 데이터베이스 연결에 실패 했어요, 재시도 중 %v... (attempt %d/%d)", delay, attempts, maxAttempts)
		time.Sleep(delay)

		if err := c.ConnectDB(); err != nil {
			log.Println("[postgre] 데이터베이스 연결에 실패 했어요:", err)
			return
		}
	}
}

func (c *Client) Dispose() error {
	postgre, err := c.DB.DB()
	if err != nil {
		return err
	}

	if err := postgre.Close(); err != nil {
		return err
	}

	return nil
}

func NewClient() *Client {
	opt := ConstructorOption{
		Phase: config.Phase(),
		DSN:   config.PostgreSQLDSN(),
	}

	client := &Client{
		opt: opt,
	}

	return client
}
