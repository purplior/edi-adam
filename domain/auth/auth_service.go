package auth

import (
	"time"

	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/lib/myjwt"
)

var (
	jwtSecretKey = []byte(config.JwtSecretKey())
)

type (
	AuthService interface {
		SignInByEmailVerification(
			ctx context.APIContext,
			request SignInByEmailVerificationRequest,
		) (
			identityToken IdentityToken,
			identity Identity,
			err error,
		)

		SignUpByEmailVerification(
			ctx context.APIContext,
			request SignUpByEmailVerificationRequest,
		) (
			err error,
		)

		RefreshIdentityToken(
			ctx context.APIContext,
			identityToken IdentityToken,
		) (
			refreshedIdentityToken IdentityToken,
			err error,
		)
	}

	authService struct {
		emailVerificationService verification.EmailVerificationService
		userService              user.UserService
	}
)

func (s *authService) SignInByEmailVerification(
	ctx context.APIContext,
	request SignInByEmailVerificationRequest,
) (
	IdentityToken,
	Identity,
	error,
) {
	existedUser, err := s.userService.GetByAccount(
		ctx,
		user.JoinMethod_Email,
		request.AccountID,
	)
	if err != nil {
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

func (s *authService) SignUpByEmailVerification(
	ctx context.APIContext,
	request SignUpByEmailVerificationRequest,
) error {
	verification, err := s.emailVerificationService.Consume(
		ctx,
		request.VerificationID,
	)
	if err != nil {
		return err
	}

	_, err = s.userService.RegisterOne(
		ctx,
		user.User{
			JoinMethod:      user.JoinMethod_Email,
			AccountID:       verification.Email,
			AccountPassword: request.Password,
			Nickname:        request.Nickname,
			Role:            user.Role_User,
		},
	)

	return err
}

func (s *authService) RefreshIdentityToken(
	ctx context.APIContext,
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

func (s *authService) makeAccessToken(
	payload map[string]interface{},
) (
	string,
	error,
) {
	// 유효 기간: 1시간
	atExpires := time.Now().Add(time.Hour).Unix()

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
	rtExpires := time.Now().Add(time.Hour * 24 * 180).Unix()

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

	at, err := s.makeAccessToken(atPayload)
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

	print(atPayload["version"])

	var identity Identity
	identity.SyncWith(atPayload)

	newAccessToken, err := s.makeAccessToken(atPayload)
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
) AuthService {
	return &authService{
		emailVerificationService: emailVerificationService,
		userService:              userService,
	}
}
