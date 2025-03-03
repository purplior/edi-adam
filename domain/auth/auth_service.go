package auth

import (
	"time"

	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/missionlog"
	"github.com/purplior/edi-adam/domain/shared/constant"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/domain/user"
	"github.com/purplior/edi-adam/domain/verification"
	"github.com/purplior/edi-adam/domain/wallet"
	"github.com/purplior/edi-adam/lib/dt"
	"github.com/purplior/edi-adam/lib/myjwt"
	"github.com/purplior/edi-adam/lib/security"
	"github.com/purplior/edi-adam/lib/strgen"
	"github.com/purplior/edi-adam/lib/validator"
)

var (
	jwtSecretKey = []byte(config.JwtSecretKey())
)

type (
	AuthService interface {
		SignIn(
			session inner.Session,
			dto SignInDTO,
		) (
			identityToken IdentityToken,
			identity inner.Identity,
			err error,
		)

		SignUp(
			session inner.Session,
			dto SignUpDTO,
		) (
			identityToken IdentityToken,
			identity inner.Identity,
			err error,
		)

		GetTempAccessToken(
			session inner.Session,
			identity inner.Identity,
		) (
			accessToken string,
			err error,
		)
	}

	authService struct {
		verificationService verification.VerificationService
		userService         user.UserService
		walletService       wallet.WalletService
		missionLogService   missionlog.MissionLogService
	}
)

func (s *authService) SignIn(
	session inner.Session,
	dto SignInDTO,
) (
	identityToken IdentityToken,
	identity inner.Identity,
	err error,
) {
	phoneNumber, _, err := s.consumeVerification(session, dto.VerificationID)
	if err != nil {
		return identityToken, identity, err
	}

	existedUser, err := s.userService.Get(
		session,
		user.QueryOption{
			PhoneNumber: phoneNumber,
		},
	)
	if err != nil {
		return identityToken, identity, exception.ErrUnauthorized
	}
	if existedUser.IsInactivated {
		return identityToken, identity, exception.ErrUnauthorized
	}

	if identityToken, identity, err = s.makeToken(existedUser); err != nil {
		return identityToken, identity, exception.ErrUnauthorized
	}

	return identityToken, identity, nil
}

func (s *authService) SignUp(
	session inner.Session,
	dto SignUpDTO,
) (
	identityToken IdentityToken,
	identity inner.Identity,
	err error,
) {
	if err = validator.CheckValidNickname(dto.Nickname); err != nil {
		return identityToken, identity, err
	}
	if err = session.BeginTransaction(); err != nil {
		return identityToken, identity, err
	}

	var phoneNumber string
	if phoneNumber, _, err = s.consumeVerification(session, dto.VerificationID); err != nil {
		session.RollbackTransaction()
		return identityToken, identity, err
	}

	// 2. 계정 생성
	me, err := s.userService.RegisterMember(
		session,
		user.RegisterMemberDTO{
			PhoneNumber: phoneNumber,
			Nickname:    dto.Nickname,
			Avatar: model.UserAvatar{
				Text:  strgen.ExtractInitialChar(dto.Nickname),
				Theme: "default",
			},
			IsMarketingAgree: dto.IsMarketingAgree,
		},
	)
	if err != nil {
		session.RollbackTransaction()
		return identityToken, identity, err
	}

	// 3. 지갑 생성
	_, err = s.walletService.Add(
		session,
		wallet.AddDTO{
			OwnerID: me.ID,
		},
	)
	if err != nil {
		session.RollbackTransaction()
		return identityToken, identity, err
	}

	// 4. 회원가입 미션 달성
	if err := s.missionLogService.Achieve(
		session,
		missionlog.AchieveDTO{
			UserID:    me.ID,
			MissionID: constant.MissionID_SignUp,
		},
	); err != nil {
		session.RollbackTransaction()
		return identityToken, identity, err
	}
	if err := s.missionLogService.Achieve(
		session,
		missionlog.AchieveDTO{
			UserID:    me.ID,
			MissionID: constant.MissionID_SignUpOpenEvent,
		},
	); err != nil {
		session.RollbackTransaction()
		return identityToken, identity, err
	}

	if err := session.CommitTransaction(); err != nil {
		return identityToken, identity, err
	}
	if identityToken, identity, err = s.makeToken(me); err != nil {
		return identityToken, identity, exception.ErrUnauthorized
	}

	return identityToken, identity, err
}

func (s *authService) GetTempAccessToken(
	session inner.Session,
	identity inner.Identity,
) (
	string,
	error,
) {
	atPayload, err := identity.ToMap()
	if err != nil {
		return "", err
	}

	atExpires := time.Now().Add(time.Hour).Unix()

	return s.makeAccessToken(atPayload, atExpires)
}

func (s *authService) consumeVerification(
	session inner.Session,
	verificationID uint,
) (
	phoneNumber string,
	m model.Verification,
	err error,
) {
	m, err = s.verificationService.Consume(
		session,
		verificationID,
	)
	if err != nil {
		return phoneNumber, m, err
	}

	var data map[string]interface{}
	data, err = security.DecryptMapDataWithAESGCM(m.Encrypted, config.SymmetricKey())
	if err != nil {
		return phoneNumber, m, err
	}
	if model.VerificationMethod(dt.Str(data["method"])) != model.VerificationMethod_Phone {
		return phoneNumber, m, exception.ErrBadRequest
	}

	phoneNumber = dt.Str(data["target"])
	if len(phoneNumber) == 0 {
		return phoneNumber, m, exception.ErrBadRequest
	}

	return phoneNumber, m, err
}

func (s *authService) makeAccessToken(
	payload map[string]interface{},
	atExpires int64,
) (
	string,
	error,
) {
	at, err := myjwt.SignWithHS256(
		payload,
		atExpires,
		jwtSecretKey,
	)
	if err != nil {
		return "", err
	}

	return at, nil
}

func (s *authService) makeToken(
	m model.User,
) (
	identityToken IdentityToken,
	identity inner.Identity,
	err error,
) {
	version := "v1"
	identity = inner.Identity{
		Version:    version,
		ID:         m.ID,
		AccountID:  m.AccountID,
		Nickname:   m.Nickname,
		Membership: "", // TODO: membership 정보 추가
		Role:       m.Role,
	}

	atExpires := time.Now().Add(time.Hour * 24 * 365).Unix()

	var atPayload map[string]interface{}
	var at string
	if atPayload, err = identity.ToMap(); err != nil {
		return identityToken, identity, err
	}
	if at, err = s.makeAccessToken(atPayload, atExpires); err != nil {
		return identityToken, identity, err
	}

	identityToken = IdentityToken{
		AccessToken: at,
	}

	return identityToken, identity, nil
}

func NewAuthService(
	verificationService verification.VerificationService,
	userService user.UserService,
	walletService wallet.WalletService,
	missionLogService missionlog.MissionLogService,
) AuthService {
	return &authService{
		verificationService: verificationService,
		userService:         userService,
		walletService:       walletService,
		missionLogService:   missionLogService,
	}
}
