package inner

import (
	"encoding/json"

	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/dt"
)

type (
	// 추후 Redis로 이동
	Identity struct {
		Version    string         `json:"version"`
		ID         uint           `json:"id"`
		AccountID  string         `json:"accountId"`
		Nickname   string         `json:"nickname"`
		Membership string         `json:"membership"`
		Role       model.UserRole `json:"role"`
	}
)

func (m *Identity) ToMap() (map[string]interface{}, error) {
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

func (m *Identity) SyncWith(data map[string]interface{}) {
	m.Version = dt.Str(data["version"])
	m.ID = dt.UInt(data["id"])
	m.AccountID = dt.Str(data["accountId"])
	m.Nickname = dt.Str(data["nickname"])
	m.Membership = dt.Str(data["membership"])
	m.Role = model.UserRole(dt.Int(data["role"]))
}
