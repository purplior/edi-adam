package verification

import (
	"fmt"
	"time"

	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/lib/mydate"
	"github.com/purplior/podoroot/lib/strgen"
	"github.com/purplior/podoroot/lib/validator"
)

type (
	PhoneVerificationService interface {
		Consume(
			ctx inner.Context,
			id string,
		) (
			PhoneVerification,
			error,
		)

		RequestCode(
			ctx inner.Context,
			phoneNumber string,
			isTestMode bool,
		) (
			PhoneVerification,
			error,
		)

		VerifyCode(
			ctx inner.Context,
			email string,
			code string,
		) (
			PhoneVerification,
			error,
		)
	}

	phoneVerificationService struct {
		phoneVerificationRepository PhoneVerificationRepository
	}
)

func (s *phoneVerificationService) Consume(
	ctx inner.Context,
	id string,
) (PhoneVerification, error) {
	verification, err := s.phoneVerificationRepository.FindOneById(ctx, id)
	if err != nil {
		return PhoneVerification{}, err
	}
	if !verification.IsVerified {
		return PhoneVerification{}, exception.ErrNotConsumed
	}
	if verification.IsConsumed {
		return PhoneVerification{}, exception.ErrAlreadyConsumed
	}

	err = s.phoneVerificationRepository.UpdateOne_isConsumed(ctx, id, true)
	if err != nil {
		return PhoneVerification{}, err
	}

	verification.IsConsumed = true

	return verification, nil
}

func (s *phoneVerificationService) RequestCode(
	ctx inner.Context,
	phoneNumber string,
	isTestMode bool,
) (PhoneVerification, error) {
	if err := validator.CheckValidPhoneNumber(phoneNumber); err != nil {
		return PhoneVerification{}, err
	}

	code := strgen.RandomNumber(6)
	verification := PhoneVerification{
		PhoneNumber: phoneNumber,
		Code:        code,
		IsVerified:  false,
		IsConsumed:  false,
		ExpiredAt:   mydate.After(time.Duration(5 * time.Minute)),
	}

	ver, err := s.phoneVerificationRepository.InsertOne(ctx, verification)
	if err != nil {
		return PhoneVerification{}, err
	}

	if !isTestMode {
		msg := "[포도쌤] 인증번호 [" + code + "] *타인에게 절대 알리지 마세요."
		// TODO: send email
		fmt.Println(msg)
	}

	return ver, nil
}

func (s *phoneVerificationService) VerifyCode(
	ctx inner.Context,
	phoneNumber string,
	code string,
) (
	PhoneVerification,
	error,
) {
	verification, err := s.phoneVerificationRepository.FindRecentOneByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return PhoneVerification{}, err
	}

	if verification.Code != code {
		return PhoneVerification{}, exception.ErrInvalidVerificationCode
	}
	if verification.IsVerified {
		return PhoneVerification{}, exception.ErrAlreadyVerified
	}

	err = s.phoneVerificationRepository.UpdateOne_IsVerified(
		ctx,
		verification.ID,
		true,
	)
	if err != nil {
		return PhoneVerification{}, nil
	}

	verification.IsVerified = true

	return verification, nil
}

func NewPhoneVerificationService(
	phoneVerificationRepository PhoneVerificationRepository,
) PhoneVerificationService {
	return &phoneVerificationService{
		phoneVerificationRepository: phoneVerificationRepository,
	}
}
