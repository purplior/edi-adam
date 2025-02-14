package database

import (
	"time"

	"github.com/purplior/sbec/infra/database/mongodb"
	"github.com/purplior/sbec/infra/database/sqldb"
)

type (
	DatabaseManager interface {
		Init() error

		Monitor() error

		Dispose() error
	}
)

type (
	databaseManager struct {
		sqlDBClient   *sqldb.Client
		mongoDBClient *mongodb.Client
	}
)

func (m *databaseManager) Init() error {
	if err := m.mongoDBClient.Connect(); err != nil {
		return err
	}
	if err := m.sqlDBClient.ConnectDB(); err != nil {
		return err
	}
	if err := m.sqlDBClient.MigrateDB(); err != nil {
		return err
	}

	return nil
}

func (m *databaseManager) Monitor() error {
	for {
		err := m.sqlDBClient.PingDB()
		if err != nil {
			// - 최대 5번의 재연결 시도
			// - 지수 백오프 시작 시간 2초
			m.sqlDBClient.ReconnectDB(5, 2*time.Second)
		}
		time.Sleep(30 * time.Second) // 30초마다 연결 상태 확인
	}
}

func (m *databaseManager) Dispose() error {
	return m.sqlDBClient.Dispose()
}

func NewDatabaseManager(
	sqlDBClient *sqldb.Client,
	mongoDBClient *mongodb.Client,
) DatabaseManager {
	return &databaseManager{
		sqlDBClient:   sqlDBClient,
		mongoDBClient: mongoDBClient,
	}
}
