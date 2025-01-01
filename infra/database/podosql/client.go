package podosql

import (
	"log"
	"os"
	"time"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/domain/shared/constant"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/entity"
	"gorm.io/driver/mysql"
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

func (c *Client) DBWithContext(ctx inner.Context) *DB {
	tx := ctx.TX(inner.TX_PodoSql)
	if tx != nil {
		return tx
	}

	return c.DB.WithContext(ctx.Value())
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

	db, err := gorm.Open(mysql.Open(c.opt.DSN), &gorm.Config{
		Logger:      dbLogger,
		PrepareStmt: true,
	})
	if err != nil {
		return err
	}

	if isDebugMode {
		log.Println("[podosql] 디버그모드가 활성화 되었어요")
		db.Debug()
	}

	c.DB = db
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)           // 유휴 상태로 유지할 수 있는 최대 연결 수
	sqlDB.SetMaxOpenConns(100)          // 데이터베이스에 열어둘 최대 연결 수
	sqlDB.SetConnMaxLifetime(time.Hour) // 연결의 최대 수명 (예: 1시간)

	if err := c.PingDB(); err != nil {
		return err
	}

	log.Println("[podosql] 데이터베이스 연결에 성공 했어요")

	return nil
}

func (c *Client) MigrateDB() error {
	// Relation 때문에 AutoMigrate 순서가 중요함
	return c.DB.AutoMigrate(
		entity.User{},
		entity.Wallet{},
		entity.Category{},
		entity.Ledger{},
		entity.Assistant{},
		entity.Mission{},
		entity.Challenge{},
		entity.CustomerVoice{},
		entity.EmailVerification{},
		entity.PhoneVerification{},
		entity.Bookmark{},
	)
}

func (c *Client) PingDB() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func (c *Client) ReconnectDB(
	maxAttempts int,
	baseDelay time.Duration,
) {
	attempts := 0
	for {
		err := c.PingDB()
		if err == nil {
			log.Println("[podosql] 데이터베이스 재연결에 성공 했어요")
			return
		}

		attempts++
		if attempts > maxAttempts {

			log.Fatalf("[podosql] 데이터베이스 재시도 연결 횟수가 최대를 초과 했어요: %v", err)
			return
		}

		// 2의 지수 증가
		delay := baseDelay * time.Duration(1<<attempts)
		log.Printf("[podosql] 데이터베이스 연결에 실패 했어요, 재시도 중 %v... (attempt %d/%d)", delay, attempts, maxAttempts)
		time.Sleep(delay)

		if err := c.ConnectDB(); err != nil {
			log.Println("[podosql] 데이터베이스 연결에 실패 했어요:", err)
			return
		}
	}
}

func (c *Client) Dispose() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}

func NewClient() *Client {
	opt := ConstructorOption{
		Phase: config.Phase(),
		DSN:   config.SqlDbDSN(),
	}

	client := &Client{
		opt: opt,
	}

	return client
}
