package inner

import (
	"context"
	"time"
)

type (
	// @MODULE {세션 생성기}
	SessionFactory interface {
		// API 레벨 컨텍스트 객체 생성
		CreateNewSession(time.Duration) (Session, context.CancelFunc)
	}
)
