package verification

import (
	"time"

	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/lib/mydate"
	"github.com/podossaem/podoroot/lib/mymail"
	"github.com/podossaem/podoroot/lib/strgen"
)

type (
	EmailVerificationService interface {
		Consume(
			ctx inner.Context,
			id string,
		) (
			EmailVerification,
			error,
		)

		RequestCode(
			ctx inner.Context,
			email string,
			isTestMode bool,
		) (
			EmailVerification,
			error,
		)

		VerifyCode(
			ctx inner.Context,
			email string,
			code string,
		) (
			EmailVerification,
			error,
		)
	}

	service struct {
		emailVerificationRepository EmailVerificationRepository
	}
)

func (s *service) Consume(
	ctx inner.Context,
	id string,
) (EmailVerification, error) {
	emailVerification, err := s.emailVerificationRepository.FindOneById(ctx, id)
	if err != nil {
		return EmailVerification{}, err
	}
	if !emailVerification.IsVerified {
		return EmailVerification{}, exception.ErrNotConsumed
	}
	if emailVerification.IsConsumed {
		return EmailVerification{}, exception.ErrAlreadyConsumed
	}

	err = s.emailVerificationRepository.UpdateOne_isConsumed(ctx, id, true)
	if err != nil {
		return EmailVerification{}, err
	}

	emailVerification.IsConsumed = true

	return emailVerification, nil
}

func (s *service) RequestCode(
	ctx inner.Context,
	email string,
	isTestMode bool,
) (EmailVerification, error) {
	code := strgen.RandomNumber(6)
	verification := EmailVerification{
		Email:      email,
		Code:       code,
		IsVerified: false,
		IsConsumed: false,
		ExpiredAt:  mydate.After(time.Duration(5 * time.Minute)),
	}

	subject := "[포도쌤] 이메일 인증요청이 왔어요"
	body := makeRequestCodeBody(code)

	ver, err := s.emailVerificationRepository.InsertOne(ctx, verification)
	if err != nil {
		return EmailVerification{}, err
	}

	if !isTestMode {
		if err := mymail.SendGmail(mymail.SendGmailRequest{
			To:           email,
			From:         config.CsEmail(),
			FromPassword: config.CsEmailPassword(),
			Subject:      subject,
			Body:         body,
		}); err != nil {
			return EmailVerification{}, err
		}
	}

	return ver, nil
}

func (s *service) VerifyCode(
	ctx inner.Context,
	email string,
	code string,
) (
	EmailVerification,
	error,
) {
	emailVerification, err := s.emailVerificationRepository.FindRecentOneByEmail(ctx, email)
	if err != nil {
		return EmailVerification{}, err
	}

	if emailVerification.Code != code {
		return EmailVerification{}, exception.ErrInvalidVerificationCode
	}
	if emailVerification.IsVerified {
		return EmailVerification{}, exception.ErrAlreadyVerified
	}

	err = s.emailVerificationRepository.UpdateOne_IsVerified(
		ctx,
		emailVerification.ID,
		true,
	)
	if err != nil {
		return EmailVerification{}, nil
	}

	emailVerification.IsVerified = true

	return emailVerification, nil
}

func makeRequestCodeBody(code string) string {
	return `<table border="0" cellpadding="0" cellspacing="0" style="min-width:600px;width:600px;max-width:600px;margin:0 auto" align="center" valign="top" role="presentation">
	<tbody>
	<tr>
	  <td dir="ltr" valign="top">
		  <div style="font-family:Roboto,'Segoe UI','Helvetica Neue',Frutiger,'Frutiger Linotype','Dejavu Sans','Trebuchet MS',Verdana,Arial,sans-serif;margin:0 auto;padding:0;max-width:600px">
				<div style="margin:0 16px 32px 16px;font-size:24px;font-weight:300">인증번호 : <span dir="ltr" style="white-space:nowrap;direction:ltr">` + code + `</span></div>
				<div style="margin:0 16px 32px 16px;font-size:24px;font-weight:300">위 코드를 입력창에 입력 해주세요.</div>
				<div style="border-top:1px solid #bdbdbd;padding-top:32px;"></div>
			</div>
		</td>
	</tr>
	</tbody>
</table>`
}

func NewEmailVerificationService(
	emailVerificationRepository EmailVerificationRepository,
) EmailVerificationService {
	return &service{
		emailVerificationRepository: emailVerificationRepository,
	}
}
