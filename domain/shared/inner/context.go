package inner

import (
	"context"

	"gorm.io/gorm"
)

const (
	TX_sqldb   TX = 1
	TX_mongodb TX = 2
)

type (
	TX int

	Context interface {
		Value() context.Context
		TX(target TX) *gorm.DB
		SetTX(target TX, tx *gorm.DB)
		ClearTX(target TX)
	}

	ContextManager interface {
		NewContext() (Context, context.CancelFunc)
		NewStreamingContext() (Context, context.CancelFunc)
		BeginTX(ctx Context, target TX) error
		CommitTX(ctx Context, target TX) error
		RollbackTX(ctx Context, target TX)
		ClearTX(ctx Context, target TX)
	}
)
