package exception

import "errors"

var (
	ErrBadRequest    = errors.New("잘못된 요청이에요")
	ErrUnauthorized  = errors.New("인증에 실패 했어요")
	ErrNotAcceptable = errors.New("받아들일 수 없는 요청이에요")
)
