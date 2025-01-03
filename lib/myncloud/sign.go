package myncloud

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

type (
	MakeSignaturePayload struct {
		HTTPMethod string
		URI        string
		AccessKey  string
		SecretKey  string
	}
)

func MakeSignature(payload MakeSignaturePayload) (string, string, error) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	message := payload.HTTPMethod + " " + payload.URI + "\n" + timestamp + "\n" + payload.AccessKey
	h := hmac.New(sha256.New, []byte(payload.SecretKey))
	_, err := h.Write([]byte(message))
	if err != nil {
		return "", timestamp, err
	}
	signingKey := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signingKey, timestamp, nil
}
