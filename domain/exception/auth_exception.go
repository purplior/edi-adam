package exception

import "errors"

var (
	ErrInvalidVerificationCode = errors.New("잘못된 인증코드 입니다")
	ErrAlreadyVerified         = errors.New("이미 인증이 완료 되었습니다")
	ErrNotConsumed             = errors.New("아직 인증되지 않았습니다")
	ErrAlreadyConsumed         = errors.New("이미 인증 되었습니다")
	ErrAlreadySignedUp         = errors.New("이미 회원가입이 완료된 계정입니다")
)
