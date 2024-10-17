package myjwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/podossaem/podoroot/domain/shared/exception"
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

func ParseWithHMACWithoutVerify(tokenString string) (map[string]interface{}, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, exception.ErrUnauthorized
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, exception.ErrUnauthorized
	}

	// 페이로드를 JSON으로 변환합니다.
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payload, &payloadMap); err != nil {
		return nil, exception.ErrUnauthorized
	}

	return payloadMap, nil
}
