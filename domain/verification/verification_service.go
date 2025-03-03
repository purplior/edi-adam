package verification

import (
	"strings"
	"time"

	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/logger"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/port/sms"
	"github.com/purplior/edi-adam/lib/dt"
	"github.com/purplior/edi-adam/lib/mydate"
	"github.com/purplior/edi-adam/lib/mymail"
	"github.com/purplior/edi-adam/lib/security"
	"github.com/purplior/edi-adam/lib/strgen"
)

type (
	VerificationService interface {
		Request(
			session inner.Session,
			method model.VerificationMethod,
			target string,
		) (
			model.Verification,
			error,
		)

		Verify(
			session inner.Session,
			id uint,
			code string,
		) error

		Consume(
			session inner.Session,
			id uint,
		) (
			model.Verification,
			error,
		)
	}

	verificationService struct {
		smsClient              *sms.Client
		verificationRepository VerificationRepository
	}
)

func (s *verificationService) Request(
	session inner.Session,
	method model.VerificationMethod,
	target string,
) (
	m model.Verification,
	err error,
) {
	switch method {
	case model.VerificationMethod_Phone:
		target = strings.ReplaceAll(target, "-", "")
	}

	logger.Debug("Request Verification Code: %s", method)

	code := strgen.RandomNumber(6)
	data := map[string]interface{}{
		"method": method,
		"target": target,
	}

	var encrypted, hash string
	if encrypted, err = security.EncryptMapDataWithAESGCM(data, config.SymmetricKey()); err != nil {
		return m, err
	}
	if hash, err = security.HashMapDataWithSHA256(data); err != nil {
		return m, err
	}

	var count int
	if count, err = s.verificationRepository.ReadCount(session, QueryOption{
		Hash:           hash,
		CreatedAtStart: mydate.DayStartFromNow(0),
		CreatedAtEnd:   mydate.DayEndFromNow(0),
	}); err != nil {
		return m, err
	}
	if count >= 5 { // 하루에 최대 5번만 요청 허용
		return m, exception.ErrPhoneVerificationExceed
	}

	logger.Debug("Today Verification Count: %d", count)

	if m, err = s.verificationRepository.Create(
		session,
		model.Verification{
			Encrypted:  encrypted,
			Hash:       hash,
			Code:       code,
			IsVerified: false,
			IsConsumed: false,
			ExpiredAt:  mydate.After(time.Duration(5 * time.Minute)),
		},
	); err != nil {
		return m, err
	}

	switch method {
	case model.VerificationMethod_Phone:
		logger.Debug("Send SMS: %s", target)
		response, err := s.smsClient.SendSMS(sms.SendSMSRequest{
			Content: "[샘비서] 인증번호 [" + code + "] *타인에게 절대 알리지 마세요.",
			ToList: []string{
				target,
			},
		})
		statusCode := dt.Int(response.StatusCode)
		logger.Debug("SMS Status Code: %d", statusCode)

		if err != nil {
			return m, err
		}
		if statusCode < 200 || statusCode >= 300 {
			return m, exception.ErrSMSFailed
		}
	case model.VerificationMethod_Email:
		if err = mymail.SendGmail(mymail.SendGmailRequest{
			To:           target,
			From:         config.CustomerVoiceEmail(),
			FromPassword: config.CustomerVoiceEmailPassword(),
			Subject:      "[샘비서] 이메일 인증요청이 왔어요",
			Body:         makeRequestEmailBody(code),
		}); err != nil {
			return m, err
		}
	}

	return m, nil
}

func (s *verificationService) Verify(
	session inner.Session,
	id uint,
	code string,
) (
	err error,
) {
	var m model.Verification
	if m, err = s.verificationRepository.Read(
		session,
		QueryOption{
			ID: id,
		},
	); err != nil {
		return err
	}

	logger.Debug("Request Verification: %d / %s", id, code)

	if m.Code != code {
		return exception.ErrInvalidVerificationCode
	}
	if m.IsVerified {
		return exception.ErrAlreadyVerified
	}

	if err = s.verificationRepository.Updates(
		session,
		QueryOption{
			ID: m.ID,
		},
		model.Verification{
			IsVerified: true,
		},
	); err != nil {
		return err
	}

	return nil
}

func (s *verificationService) Consume(
	session inner.Session,
	id uint,
) (
	m model.Verification,
	err error,
) {
	m, err = s.verificationRepository.Read(
		session,
		QueryOption{
			ID: id,
		},
	)
	if err != nil {
		return m, err
	}
	if !m.IsVerified {
		return m, exception.ErrNotConsumed
	}
	if m.IsConsumed {
		return m, exception.ErrAlreadyConsumed
	}

	err = s.verificationRepository.Updates(
		session,
		QueryOption{
			ID: id,
		},
		model.Verification{
			IsConsumed: true,
		},
	)
	if err != nil {
		return m, err
	}

	return m, nil
}

func makeRequestEmailBody(code string) string {
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

func NewVerificationService(
	smsClient *sms.Client,
	verificationRepository VerificationRepository,
) VerificationService {
	return &verificationService{
		smsClient:              smsClient,
		verificationRepository: verificationRepository,
	}
}
