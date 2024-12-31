package exception

import "errors"

var (
	ErrInvalidVerificationCode = errors.New("잘못된 인증코드에요")
	ErrAlreadyVerified         = errors.New("이미 인증이 완료 되었어요")
	ErrNotConsumed             = errors.New("아직 인증되지 않았어요")
	ErrAlreadyConsumed         = errors.New("이미 인증 되었어요")
	ErrAlreadySignedUp         = errors.New("이미 회원가입이 완료된 계정이에요")
	ErrNoSignedUpPhone         = errors.New("아직 가입하지 않은 휴대번호에요")
	ErrPhoneVerificationExceed = errors.New("휴대폰 인증의 하루 허용횟수(5회)를 초과했어요")
	ErrNotAllowedNickname      = errors.New("허용하지 않는 닉네임이에요")
)
