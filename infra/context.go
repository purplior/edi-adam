package infra

import (
	"context"
	"time"

	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/infra/database/podopaysql"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"gorm.io/gorm"
)

type (
	ctx struct {
		value        context.Context
		podosqlTX    *gorm.DB
		podopaysqlTX *gorm.DB
	}
)

func (c *ctx) Value() context.Context {
	return c.value
}

func (c *ctx) TX(target inner.TX) *gorm.DB {
	switch target {
	case inner.TX_PodoSql:
		return c.podosqlTX
	case inner.TX_PodopaySql:
		return c.podopaysqlTX
	}

	return nil
}

func (c *ctx) SetTX(target inner.TX, tx *gorm.DB) {
	switch target {
	case inner.TX_PodoSql:
		c.podosqlTX = tx
	case inner.TX_PodopaySql:
		c.podopaysqlTX = tx
	}
}

func (c *ctx) ClearTX(target inner.TX) {
	switch target {
	case inner.TX_PodoSql:
		c.podosqlTX = nil
	case inner.TX_PodopaySql:
		c.podopaysqlTX = nil
	}
}

type (
	contextManager struct {
		podosqlClient    *podosql.Client
		podopaysqlClient *podopaysql.Client
	}
)

func (c *contextManager) NewContext() (inner.Context, context.CancelFunc) {
	todoCtx := context.TODO()
	value, cancel := context.WithTimeout(todoCtx, time.Duration(12*time.Second))

	return &ctx{
		value:        value,
		podosqlTX:    nil,
		podopaysqlTX: nil,
	}, cancel
}

func (c *contextManager) NewStreamingContext() (inner.Context, context.CancelFunc) {
	todoCtx := context.TODO()
	value, cancel := context.WithTimeout(todoCtx, time.Duration(5*time.Minute))

	return &ctx{
		value:        value,
		podosqlTX:    nil,
		podopaysqlTX: nil,
	}, cancel
}

func (c *contextManager) BeginTX(ctx inner.Context, target inner.TX) error {
	tx := ctx.TX(target)
	if tx != nil {
		return exception.ErrInTransaction
	}

	switch target {
	case inner.TX_PodoSql:
		tx := c.podosqlClient.WithContext(ctx.Value()).Begin()
		ctx.SetTX(target, tx)
	case inner.TX_PodopaySql:
		tx := c.podopaysqlClient.WithContext(ctx.Value()).Begin()
		ctx.SetTX(target, tx)
	}

	return nil
}

func (c *contextManager) CommitTX(ctx inner.Context, target inner.TX) error {
	tx := ctx.TX(target)
	if tx == nil {
		return exception.ErrInvalidTransaction
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	ctx.ClearTX(target)

	return err
}

func (c *contextManager) RollbackTX(ctx inner.Context, target inner.TX) {
	tx := ctx.TX(target)
	if tx != nil {
		tx.Rollback()
	}
	ctx.ClearTX(target)
}

func NewContextManager(
	podosqlClient *podosql.Client,
	podopaysqlClient *podopaysql.Client,
) inner.ContextManager {
	return &contextManager{
		podosqlClient:    podosqlClient,
		podopaysqlClient: podopaysqlClient,
	}
}
