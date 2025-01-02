package exception

import "errors"

var (
	ErrInternalServer  = errors.New("일시적인 서버오류가 발생했어요")
	ErrTest            = errors.New("테스트용 에러에요")
	ErrBadRequest      = errors.New("잘못된 요청이에요")
	ErrUnauthorized    = errors.New("인증에 실패 했어요")
	ErrNotAcceptable   = errors.New("받아들일 수 없는 요청이에요")
	ErrNotFound        = errors.New("리소스를 찾을 수 없어요")
	ErrNoPodo          = errors.New("사용할 수 있는 포도가 더 필요해요")
	ErrAlreadyReceived = errors.New("이미 수령했어요")
)
