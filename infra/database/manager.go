package database

import (
	"time"

	"github.com/podossaem/podoroot/infra/database/podosql"
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
		podosqlClient *podosql.Client
	}
)

func (m *databaseManager) Init() error {
	if err := m.podosqlClient.ConnectDB(); err != nil {
		return err
	}
	if err := m.podosqlClient.MigrateDB(); err != nil {
		return err
	}

	return nil
}

func (m *databaseManager) Monitor() error {
	for {
		err := m.podosqlClient.PingDB()
		if err != nil {
			// - 최대 5번의 재연결 시도
			// - 지수 백오프 시작 시간 2초
			m.podosqlClient.ReconnectDB(5, 2*time.Second)
		}
		time.Sleep(30 * time.Second) // 30초마다 연결 상태 확인
	}
}

func (m *databaseManager) Dispose() error {
	return m.podosqlClient.Dispose()
}

func NewDatabaseManager(
	podosqlClient *podosql.Client,
) DatabaseManager {
	return &databaseManager{
		podosqlClient: podosqlClient,
	}
}
