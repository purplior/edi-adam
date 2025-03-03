package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

func HashMapDataWithSHA256(data map[string]interface{}) (string, error) {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(plaintext)

	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

func SignMapDataWithHMACSHA256(data map[string]interface{}, secretKey string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(jsonData)
	return hex.EncodeToString(h.Sum(nil)), nil
}
