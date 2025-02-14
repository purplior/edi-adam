package sms

import "errors"

var (
	ErrInvalidContentSize = errors.New("컨텐츠의 사이즈가 유효하지 않아요")
	ErrSendSMS            = errors.New("문자 발송 중 오류가 발생했어요")
)
