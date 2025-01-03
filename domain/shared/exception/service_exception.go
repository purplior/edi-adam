package exception

import "errors"

var (
	ErrSMSFailed = errors.New("SMS 발송에 실패했어요")
)
