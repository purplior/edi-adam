package session

import (
	"context"
	"time"

	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	sessionFactory struct {
		postgreClient *postgre.Client
	}
)

func (sf *sessionFactory) CreateNewSession(time.Duration) (inner.Session, context.CancelFunc) {
	_ctx := context.TODO()
	ctx, cancel := context.WithTimeout(_ctx, time.Duration(12*time.Second))

	return &_session{
		identity:      nil,
		innerContext:  ctx,
		postgreTX:     nil,
		postgreClient: sf.postgreClient,
	}, cancel
}

func NewFactory(
	postgreClient *postgre.Client,
) inner.SessionFactory {
	return &sessionFactory{
		postgreClient: postgreClient,
	}
}
