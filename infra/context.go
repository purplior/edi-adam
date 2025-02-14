package infra

import (
	"context"
	"time"

	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/infra/database/sqldb"
	"gorm.io/gorm"
)

type (
	ctx struct {
		value   context.Context
		sqldbTX *gorm.DB
	}
)

func (c *ctx) Value() context.Context {
	return c.value
}

func (c *ctx) TX(target inner.TX) *gorm.DB {
	switch target {
	case inner.TX_sqldb:
		return c.sqldbTX
	}

	return nil
}

func (c *ctx) SetTX(target inner.TX, tx *gorm.DB) {
	switch target {
	case inner.TX_sqldb:
		c.sqldbTX = tx
	}
}

func (c *ctx) ClearTX(target inner.TX) {
	switch target {
	case inner.TX_sqldb:
		c.sqldbTX = nil
	}
}

type (
	contextManager struct {
		sqldbClient *sqldb.Client
	}
)

func (c *contextManager) NewContext() (inner.Context, context.CancelFunc) {
	todoCtx := context.TODO()
	value, cancel := context.WithTimeout(todoCtx, time.Duration(12*time.Second))

	return &ctx{
		value:   value,
		sqldbTX: nil,
	}, cancel
}

func (c *contextManager) NewStreamingContext() (inner.Context, context.CancelFunc) {
	todoCtx := context.TODO()
	value, cancel := context.WithTimeout(todoCtx, time.Duration(5*time.Minute))

	return &ctx{
		value:   value,
		sqldbTX: nil,
	}, cancel
}

func (c *contextManager) BeginTX(ctx inner.Context, target inner.TX) error {
	tx := ctx.TX(target)
	if tx != nil {
		return exception.ErrInTransaction
	}

	switch target {
	case inner.TX_sqldb:
		tx := c.sqldbClient.WithContext(ctx.Value()).Begin()
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

func (c *contextManager) ClearTX(ctx inner.Context, target inner.TX) {
	ctx.ClearTX(target)
}

func NewContextManager(
	sqldbClient *sqldb.Client,
) inner.ContextManager {
	return &contextManager{
		sqldbClient: sqldbClient,
	}
}
