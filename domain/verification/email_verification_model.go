package verification

import "time"

type (
	EmailVerification struct {
		ID         string    `json:"id"`
		Email      string    `json:"email"`
		Code       string    `json:"code"`
		IsVerified bool      `json:"verified"`
		IsConsumed bool      `json:"consumed"`
		ExpiredAt  time.Time `json:"expiredAt"`
		CreatedAt  time.Time `json:"createdAt"`
	}
)
