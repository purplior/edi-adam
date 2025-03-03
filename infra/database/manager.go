package database

import (
	"time"

	"github.com/purplior/edi-adam/infra/database/dynamo"
	"github.com/purplior/edi-adam/infra/database/postgre"
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
		postgreClient *postgre.Client
		dyanmoClient  *dynamo.Client
	}
)

func (m *databaseManager) Init() error {
	if err := m.postgreClient.ConnectDB(); err != nil {
		return err
	}
	if err := m.postgreClient.MigrateDB(); err != nil {
		return err
	}

	return nil
}

func (m *databaseManager) Monitor() error {
	for {
		err := m.postgreClient.PingDB()
		if err != nil {
			// - 최대 5번의 재연결 시도
			// - 지수 백오프 시작 시간 2초
			m.postgreClient.ReconnectDB(5, 2*time.Second)
		}
		time.Sleep(30 * time.Second) // 30초마다 연결 상태 확인
	}
}

func (m *databaseManager) Dispose() error {
	return m.postgreClient.Dispose()
}

func NewDatabaseManager(
	postgreClient *postgre.Client,
	dynamoClient *dynamo.Client,
) DatabaseManager {
	return &databaseManager{
		postgreClient: postgreClient,
		dyanmoClient:  dynamoClient,
	}
}
