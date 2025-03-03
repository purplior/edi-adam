package model

import "time"

type (
	VerificationMethod string
)

const (
	VerificationMethod_Phone VerificationMethod = "phone"
	VerificationMethod_Email VerificationMethod = "email"
)

type (
	// @MODEL {인증 (이메일/휴대폰)}
	Verification struct {
		ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		// 암호화 되어 저장됨
		// { "method": "email", "target": "abc@purplior.com" }
		// { "method": "phone", "target": "01012345678" }
		Encrypted string `gorm:"size:100" json:"email"`
		Hash      string `gorm:"size:80" json:"hash"` // 동일한 요청이 몇개 왔는지 확인하기 위한 식별자 역할
		Code      string `gorm:"size:10" json:"code"`

		IsVerified bool `json:"isVerified"` // 검증이 완료된 여부
		IsConsumed bool `json:"isConsumed"` // 검증이 완료되고 소비될 수 있는데, 한번 소비된 인증은 다시 소비하지 못함

		ExpiredAt time.Time `json:"expiredAt"`
		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	}
)
