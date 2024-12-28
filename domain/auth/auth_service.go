package auth

import (
	"time"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/domain/challenge"
	"github.com/purplior/podoroot/domain/shared/constant"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/user"
	"github.com/purplior/podoroot/domain/verification"
	"github.com/purplior/podoroot/domain/wallet"
	"github.com/purplior/podoroot/lib/myjwt"
	"github.com/purplior/podoroot/lib/strgen"
)

var (
	jwtSecretKey = []byte(config.JwtSecretKey())
)

type (
	AuthService interface {
		SignIn_ByEmailVerification(
			ctx inner.Context,
			request SignInByEmailVerificationRequest,
		) (
			identityToken IdentityToken,
			identity Identity,
			err error,
		)

		SignUp_ByEmailVerification(
			ctx inner.Context,
			request SignUpByEmailVerificationRequest,
		) (
			err error,
		)

		RefreshIdentityToken(
			ctx inner.Context,
			identityToken IdentityToken,
		) (
			refreshedIdentityToken IdentityToken,
			err error,
		)

		GetTempAccessToken(
			ctx inner.Context,
			identity Identity,
		) (
			accessToken string,
			err error,
		)
	}

	authService struct {
		emailVerificationService verification.EmailVerificationService
		userService              user.UserService
		walletService            wallet.WalletService
		challengeService         challenge.ChallengeService
		cm                       inner.ContextManager
	}
)

func (s *authService) SignIn_ByEmailVerification(
	ctx inner.Context,
	request SignInByEmailVerificationRequest,
) (
	IdentityToken,
	Identity,
	error,
) {
	existedUser, err := s.userService.GetOne_ByAccount(
		ctx,
		user.JoinMethod_Email,
		request.AccountID,
	)
	if err != nil {
		return IdentityToken{}, Identity{}, exception.ErrUnauthorized
	}
	if existedUser.IsInactivated {
		return IdentityToken{}, Identity{}, exception.ErrUnauthorized
	}
	if err := existedUser.ComparePassword(request.Password); err != nil {
		return IdentityToken{}, Identity{}, exception.ErrUnauthorized
	}

	identityToken, identity, err := s.makeToken(existedUser)
	if err != nil {
		return IdentityToken{}, Identity{}, exception.ErrUnauthorized
	}

	return identityToken, identity, nil
}

func (s *authService) SignUp_ByEmailVerification(
	ctx inner.Context,
	request SignUpByEmailVerificationRequest,
) error {
	if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
			panic(r)
		}
	}()

	// 1. 이메일 검증
	verification, err := s.emailVerificationService.Consume(
		ctx,
		request.VerificationID,
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	// 2. 계정 생성
	me, err := s.userService.RegisterOne(
		ctx,
		user.User{
			JoinMethod:       user.JoinMethod_Email,
			AccountID:        verification.Email,
			AccountPassword:  request.Password,
			AvatarTheme:      1,
			AvatarText:       strgen.ExtractInitialChar(request.Nickname),
			Nickname:         request.Nickname,
			Role:             user.Role_User,
			IsMarketingAgree: request.IsMarketingAgree,
		},
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	// 3. 지갑 생성
	_, err = s.walletService.RegisterOne(
		ctx,
		wallet.Wallet{
			OwnerID: me.ID,
			Podo:    0,
		},
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	// 4. 회원가입 미션 달성
	if err := s.challengeService.AchieveOne_ByUserAndMission(
		ctx,
		me.ID,
		constant.MissionID_SignUp,
	); err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	return nil
}

func (s *authService) RefreshIdentityToken(
	ctx inner.Context,
	identityToken IdentityToken,
) (
	IdentityToken,
	error,
) {
	rtPayload, err := s.getRefreshTokenPayload(identityToken.RefreshToken)
	if err != nil {
		return IdentityToken{}, exception.ErrUnauthorized
	}

	identity, newAt, err := s.getIdentityAndNewAccessTokenWithoutVerify(identityToken.AccessToken)
	if err != nil {
		return IdentityToken{}, exception.ErrUnauthorized
	}

	if rtPayload.ID != identity.ID || rtPayload.Version != identity.Version {
		return IdentityToken{}, exception.ErrUnauthorized
	}

	return IdentityToken{
		AccessToken:  newAt,
		RefreshToken: identityToken.RefreshToken,
	}, nil
}

func (s *authService) GetTempAccessToken(
	ctx inner.Context,
	identity Identity,
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

func (s *authService) makeAccessToken(
	payload map[string]interface{},
	atExpires int64,
) (
	string,
	error,
) {
	// 유효 기간: 1년
	// atExpires := time.Now().Add(time.Hour * 24 * 365).Unix()

	// 임시 10초
	// atExpires := time.Now().Add(time.Second * 10).Unix()
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

func (s *authService) makeRefreshToken(
	payload map[string]interface{},
) (
	string,
	error,
) {
	// 유효 기간: 6개월
	rtExpires := time.Now().Add(time.Hour * 24 * 365).Unix()

	// 유효 기간: 1분
	// rtExpires := time.Now().Add(time.Minute).Unix()
	rt, err := myjwt.SignWithHS256(
		payload,
		rtExpires,
		jwtSecretKey,
	)
	if err != nil {
		return "", err
	}

	return rt, nil
}

func (s *authService) makeToken(
	user user.User,
) (
	IdentityToken,
	Identity,
	error,
) {
	version := "v1"
	identity := Identity{
		Version:    version,
		ID:         user.ID,
		JoinMethod: user.JoinMethod,
		AccountID:  user.AccountID,
		Nickname:   user.Nickname,
		Role:       user.Role,
	}

	atPayload, err := identity.ToMap()
	if err != nil {
		return IdentityToken{}, Identity{}, err
	}

	atExpires := time.Now().Add(time.Hour * 24 * 365).Unix()
	at, err := s.makeAccessToken(atPayload, atExpires)
	if err != nil {
		return IdentityToken{}, Identity{}, err
	}

	refreshTokenPayload := RefreshTokenPayload{
		Version: version,
		ID:      user.ID,
	}
	rtPayload, err := refreshTokenPayload.ToMap()
	if err != nil {
		return IdentityToken{}, Identity{}, err
	}

	rt, err := s.makeRefreshToken(rtPayload)
	if err != nil {
		return IdentityToken{}, Identity{}, err
	}

	return IdentityToken{
		AccessToken:  at,
		RefreshToken: rt,
	}, identity, nil
}

func (s *authService) getIdentityAndNewAccessTokenWithoutVerify(
	accessToken string,
) (
	Identity,
	string,
	error,
) {
	atPayload, _ := myjwt.ParseWithHMACWithoutVerify(accessToken)

	var identity Identity
	identity.SyncWith(atPayload)

	// 유효 기간: 1년
	atExpires := time.Now().Add(time.Hour * 24 * 365).Unix()
	newAccessToken, err := s.makeAccessToken(atPayload, atExpires)
	if err != nil {
		return Identity{}, "", err
	}

	return identity, newAccessToken, nil
}

func (s *authService) getRefreshTokenPayload(
	refreshToken string,
) (
	RefreshTokenPayload,
	error,
) {
	payload, err := myjwt.ParseWithHMAC(refreshToken, jwtSecretKey)
	if err != nil {
		return RefreshTokenPayload{}, err
	}

	var rtPayload RefreshTokenPayload
	rtPayload.SyncWith(payload)

	return rtPayload, nil
}

func NewAuthService(
	emailVerificationService verification.EmailVerificationService,
	userService user.UserService,
	walletService wallet.WalletService,
	challengeService challenge.ChallengeService,
	cm inner.ContextManager,
) AuthService {
	return &authService{
		emailVerificationService: emailVerificationService,
		userService:              userService,
		walletService:            walletService,
		challengeService:         challengeService,
		cm:                       cm,
	}
}
