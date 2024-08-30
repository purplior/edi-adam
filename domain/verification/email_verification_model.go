package verification

import "time"

type (
	EmailVerification struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		Code      string    `json:"code"`
		Verified  bool      `json:"verified"`
		Consumed  bool      `json:"consumed"`
		Ignored   bool      `json:"ignored"`
		ExpiredAt time.Time `json:"expiredAt"`
		CreatedAt time.Time `json:"createdAt"`
	}
)
