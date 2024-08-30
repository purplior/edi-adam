package context

import (
	"context"
	"time"
)

type (
	APIContext context.Context
)

func NewAPIContext() (context.Context, context.CancelFunc) {
	ctx := context.TODO()

	return context.WithTimeout(ctx, time.Duration(12*time.Second))
}
