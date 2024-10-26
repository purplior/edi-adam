package podopaysql

import (
	"context"
	"log"
	"time"

	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/shared/constant"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/infra/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	ConstructorOption struct {
		Phase constant.Phase
		DSN   string
	}

	Client struct {
		*gorm.DB
		opt ConstructorOption
		tx  *gorm.DB
	}
)

func (c *Client) DBWithContext(ctx context.Context) *gorm.DB {
	if c.tx != nil {
		return c.tx
	}

	return c.DB.WithContext(ctx)
}

func (c *Client) ConnectDB() error {
	db, err := gorm.Open(postgres.Open(c.opt.DSN), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
	})
	if err != nil {
		return err
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

	log.Println("[podopaysql] 데이터베이스 연결에 성공 했어요")

	return nil
}

func (c *Client) MigrateDB() error {
	return c.DB.AutoMigrate(
		entity.Wallet{},
		entity.Ledger{},
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
			log.Println("[podopaysql] 데이터베이스 재연결에 성공 했어요")
			return
		}

		attempts++
		if attempts > maxAttempts {

			log.Fatalf("[podopaysql] 데이터베이스 재시도 연결 횟수가 최대를 초과 했어요: %v", err)
			return
		}

		// 2의 지수 증가
		delay := baseDelay * time.Duration(1<<attempts)
		log.Printf("[podopaysql] 데이터베이스 연결에 실패 했어요, 재시도 중 %v... (attempt %d/%d)", delay, attempts, maxAttempts)
		time.Sleep(delay)

		if err := c.ConnectDB(); err != nil {
			log.Println("[podopaysql] 데이터베이스 연결에 실패 했어요:", err)
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

func (c *Client) BeginTX(ctx context.Context) error {
	if c.tx != nil {
		return exception.ErrInTransaction
	}

	c.tx = c.WithContext(ctx).Begin()

	return nil
}

func (c *Client) CommitTX() error {
	if c.tx == nil {
		return exception.ErrNoTransaction
	}

	err := c.tx.Commit().Error
	c.tx = nil

	return err
}

func (c *Client) RollbackTX() {
	if c.tx != nil {
		c.tx.Rollback()
		c.tx = nil
	}
}

func (c *Client) RecoverTX() {
	if c.tx != nil {
		if r := recover(); r != nil {
			c.tx.Rollback()
			c.tx = nil
		} else {
			c.tx = nil
		}
	}
}

func NewClient() *Client {
	opt := ConstructorOption{
		Phase: config.Phase(),
		DSN:   config.PodopaySqlDSN(),
	}

	client := &Client{
		opt: opt,
		tx:  nil,
	}

	return client
}
