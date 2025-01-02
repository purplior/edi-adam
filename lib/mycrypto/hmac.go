package mycrypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func SignMapDataWithHMACSHA256(data map[string]interface{}, secretKey string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(jsonData)
	return hex.EncodeToString(h.Sum(nil)), nil
}
