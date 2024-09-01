package myjwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/podossaem/podoroot/domain/exception"
)

func SignWithHS256(
	data map[string]interface{},
	exp int64,
	secret []byte,
) (string, error) {
	claims := jwt.MapClaims{}
	for field, val := range data {
		claims[field] = val
	}
	claims["exp"] = exp

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	return tokenString, err
}

func ParseWithHMAC(tokenString string, secret []byte) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, exception.ErrUnauthorized
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, exception.ErrUnauthorized
}
