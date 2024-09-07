package auth

import (
	"time"

	"github.com/podossaem/podoroot/application/config"
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/domain/exception"
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
			err error,
		)

		SignUpByEmailVerification(
			ctx context.APIContext,
			request SignUpByEmailVerificationRequest,
		) (
			identityToken IdentityToken,
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
	error,
) {
	existedUser, err := s.userService.GetByAccount(
		ctx,
		user.JoinMethod_Email,
		request.AccountID,
	)
	if err != nil {
		return IdentityToken{}, exception.ErrUnauthorized
	}

	if err := existedUser.ComparePassword(request.Password); err != nil {
		return IdentityToken{}, exception.ErrUnauthorized
	}

	identityToken, err := s.makeToken(existedUser)
	if err != nil {
		return IdentityToken{}, exception.ErrUnauthorized
	}

	return identityToken, nil
}

func (s *authService) SignUpByEmailVerification(
	ctx context.APIContext,
	request SignUpByEmailVerificationRequest,
) (
	IdentityToken,
	error,
) {
	verification, err := s.emailVerificationService.Consume(
		ctx,
		request.VerificationID,
	)
	if err != nil {
		return IdentityToken{}, err
	}

	newUser, err := s.userService.CreateOne(
		ctx,
		user.User{
			JoinMethod:      user.JoinMethod_Email,
			AccountID:       verification.Email,
			AccountPassword: request.Password,
			Nickname:        request.Nickname,
			Role:            user.Role_User,
		},
	)
	if err != nil {
		return IdentityToken{}, err
	}

	identityToken, err := s.makeToken(newUser)
	if err != nil {
		return IdentityToken{}, err
	}

	return identityToken, nil
}

func (s *authService) makeToken(
	user user.User,
) (
	IdentityToken,
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
		return IdentityToken{}, err
	}

	// 유효 기간: 1시간
	// atExpires := time.Now().Add(time.Hour).Unix()

	// 임시 10초
	atExpires := time.Now().Add(time.Second * 10).Unix()
	at, err := myjwt.SignWithHS256(
		atPayload,
		atExpires,
		jwtSecretKey,
	)
	if err != nil {
		return IdentityToken{}, err
	}

	refreshTokenPayload := RefreshTokenPayload{
		Version: version,
		ID:      user.ID,
	}
	rtPayload, err := refreshTokenPayload.ToMap()
	if err != nil {
		return IdentityToken{}, err
	}

	// 유효 기간: 6개월
	rtExpires := time.Now().Add(time.Hour * 24 * 180).Unix()
	rt, err := myjwt.SignWithHS256(
		rtPayload,
		rtExpires,
		jwtSecretKey,
	)
	if err != nil {
		return IdentityToken{}, err
	}

	return IdentityToken{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
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
