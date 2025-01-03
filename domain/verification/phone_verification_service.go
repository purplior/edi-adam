package verification

import (
	"strings"
	"time"

	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/port/podosms"
	"github.com/purplior/podoroot/lib/dt"
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
		smsClient                   *podosms.Client
		phoneVerificationRepository PhoneVerificationRepository
		maxCount                    int
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

	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	count, err := s.phoneVerificationRepository.FindCount_ByPhoneNumber(
		ctx,
		phoneNumber,
	)
	if err != nil {
		return PhoneVerification{}, err
	}
	if count >= s.maxCount {
		return PhoneVerification{}, exception.ErrPhoneVerificationExceed
	}

	code := strgen.RandomNumber(6)
	var requestId string = ""
	if !isTestMode {
		messageContent := "[포도쌤] 인증번호 [" + code + "] *타인에게 절대 알리지 마세요."
		response, err := s.smsClient.SendSMS(podosms.SendSMSRequest{
			Content: messageContent,
			ToList: []string{
				phoneNumber,
			},
		})
		if err != nil {
			return PhoneVerification{}, err
		}
		statusCode := dt.Int(response.StatusCode)
		if statusCode < 200 || statusCode >= 300 {
			return PhoneVerification{}, exception.ErrSMSFailed
		}

		requestId = response.RequestID
	}

	verification := PhoneVerification{
		PhoneNumber: phoneNumber,
		Code:        code,
		ReferenceID: requestId,
		IsVerified:  false,
		IsConsumed:  false,
		ExpiredAt:   mydate.After(time.Duration(5 * time.Minute)),
	}

	ver, err := s.phoneVerificationRepository.InsertOne(ctx, verification)
	if err != nil {
		return PhoneVerification{}, err
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
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	verification, err := s.phoneVerificationRepository.FindRecentOne_ByPhoneNumber(ctx, phoneNumber)
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
	smsClient *podosms.Client,
	phoneVerificationRepository PhoneVerificationRepository,
) PhoneVerificationService {
	return &phoneVerificationService{
		smsClient:                   smsClient,
		phoneVerificationRepository: phoneVerificationRepository,
		maxCount:                    5,
	}
}
