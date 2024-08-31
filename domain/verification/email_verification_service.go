package verification

import (
	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/lib/mymail"
	"github.com/podossaem/podoroot/lib/strgen"
)

type (
	EmailVerificationService interface {
		RequestCode(
			ctx context.APIContext,
			email string,
		) (
			EmailVerification,
			error,
		)

		VerifyCode(
			ctx context.APIContext,
			email string,
			code string,
		) (
			EmailVerification,
			error,
		)
	}

	service struct {
		repository EmailVerificationRepository
	}
)

func (s *service) RequestCode(
	ctx context.APIContext,
	email string,
) (EmailVerification, error) {
	code := strgen.RandomNumber(6)
	verification := EmailVerification{
		Email:      email,
		Code:       code,
		IsConsumed: false,
		IsVerified: false,
	}

	subject := "[포도쌤] 이메일 인증"
	body := makeRequestCodeBody(code)

	if err := mymail.SendGmail(mymail.SendGmailRequest{
		To:           email,
		From:         config.CsEmail(),
		FromPassword: config.CsEmailPassword(),
		Subject:      subject,
		Body:         body,
	}); err != nil {
		return EmailVerification{}, err
	}

	return s.repository.InsertOne(ctx, verification)
}

func (s *service) VerifyCode(
	ctx context.APIContext,
	email string,
	code string,
) (
	EmailVerification,
	error,
) {
	emailVerification, err := s.repository.FindOneByEmail(ctx, email)
	if err != nil {
		return EmailVerification{}, err
	}
	if emailVerification.Code != code {
		return EmailVerification{}, ErrInvalidCode
	}
	if emailVerification.IsVerified {
		return EmailVerification{}, ErrAlreadyVerified
	}

	return s.repository.UpdateOne_IsVerified(
		ctx,
		emailVerification.ID,
		true,
	)
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
	repository EmailVerificationRepository,
) EmailVerificationService {
	return &service{
		repository: repository,
	}
}
