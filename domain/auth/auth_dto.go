package auth

import (
	"encoding/json"

	"github.com/purplior/edi-adam/lib/dt"
)

type (
	RefreshTokenPayload struct {
		Version string `json:"version"`
		ID      string `json:"id"`
	}
)

func (m *RefreshTokenPayload) ToMap() (map[string]interface{}, error) {
	var data map[string]interface{}
	record, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(record, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *RefreshTokenPayload) SyncWith(data map[string]interface{}) {
	m.Version = dt.Str(data["version"])
	m.ID = dt.Str(data["id"])
}

type (
	IdentityToken struct {
		AccessToken string `json:"accessToken"`
	}

	SignInDTO struct {
		VerificationID uint `json:"verificationId"`
	}

	SignUpDTO struct {
		VerificationID   uint   `json:"verificationId"`
		Nickname         string `json:"nickname"`
		IsMarketingAgree bool   `json:"isMarketingAgree"`
	}
)
