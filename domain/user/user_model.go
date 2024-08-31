package user

import "time"

type (
	User struct {
		ID              string    `json:"id"`
		JoinMethod      string    `json:"joinMethod"`
		AccountID       string    `json:"accountId"`
		AccountPassword string    `json:"accountPassword"`
		Nickname        string    `json:"nickname"`
		Role            int       `json:"role"`
		CreatedAt       time.Time `json:"createdAt"`
	}
)

const (
	SignUpMethod_Email = "email"
	Role_Normal        = 100
	Role_Master        = 10000
)
