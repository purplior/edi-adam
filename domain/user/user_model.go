package user

import "time"

type (
	User struct {
		ID              string    `json:"id,omitempty"`
		JoinMethod      string    `json:"joinMethod"`
		AccountID       string    `json:"accountId"`
		AccountPassword string    `json:"accountPassword"`
		Nickname        string    `json:"nickname"`
		Role            int       `json:"role"`
		CreatedAt       time.Time `json:"createdAt"`
	}

	SignUpRequest struct {
		VerificationID string `json:"verificationId"`
		Password       string `json:"password"`
		Nickname       string `json:"nickname"`
	}
)

const (
	JoinMethod_Email = "email"
	Role_User        = 100
	Role_Master      = 10000
)
