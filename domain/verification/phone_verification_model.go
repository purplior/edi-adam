package verification

import "time"

type (
	PhoneVerification struct {
		ID          string    `json:"id"`
		PhoneNumber string    `json:"phoneNumber"`
		Code        string    `json:"code"`
		IsVerified  bool      `json:"verified"`
		IsConsumed  bool      `json:"consumed"`
		ExpiredAt   time.Time `json:"expiredAt"`
		CreatedAt   time.Time `json:"createdAt"`
	}
)
