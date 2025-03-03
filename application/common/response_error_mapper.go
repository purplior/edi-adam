package common

import (
	"github.com/purplior/edi-adam/application/response"
	"github.com/purplior/edi-adam/domain/shared/exception"
)

// 에러가 던져질 경우, 500응답이 아닌 적절한 상태 코드로 매핑해준다
// - 500 상태코드가 아닌 경우 에러의 메세지 텍스트 값을 바로 보내준다.
func getResponseOfError(err error) (
	statusCode int,
	message string,
) {
	statusCode = response.Status_InternalServerError
	message = "일시적인 서버 오류가 발생했어요"

	switch err {
	case exception.ErrBadRequest:
		statusCode = response.Status_BadRequest
	case exception.ErrAlreadyConsumed:
		statusCode = response.Status_BadRequest
	case exception.ErrNotConsumed:
		statusCode = response.Status_BadRequest
	case exception.ErrInvalidVerificationCode:
		statusCode = response.Status_BadRequest
	case exception.ErrAlreadyVerified:
		statusCode = response.Status_BadRequest
	case exception.ErrAlreadyReceived:
		statusCode = response.Status_BadRequest
	case exception.ErrNotAcceptable:
		statusCode = response.Status_NotAcceptable
	case exception.ErrAlreadySignedUp:
		statusCode = response.Status_NotAcceptable
	case exception.ErrUnauthorized:
		statusCode = response.Status_Unauthorized
	case exception.ErrNoRecord:
		statusCode = response.Status_NotFound
	case exception.ErrNoSignedUpPhone:
		statusCode = response.Status_Unprocessable
	case exception.ErrPhoneVerificationExceed:
		statusCode = response.Status_Unprocessable
	case exception.ErrNotAllowedNickname:
		statusCode = response.Status_BadRequest
	case exception.ErrNoCoin:
		statusCode = response.Status_Forbidden
	}
	if statusCode != response.Status_InternalServerError {
		message = err.Error()
	}

	return statusCode, message
}
