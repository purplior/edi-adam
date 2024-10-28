package inner

import (
	"context"

	"gorm.io/gorm"
)

const (
	TX_PodoSql    TX = 1
	TX_PodopaySql TX = 2
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
	}
)
