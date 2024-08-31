package credential

import (
	"encoding/json"

	"github.com/podossaem/podoroot/lib/dt"
)

type (
	Identity struct {
		Version   string `json:"version"`
		AccountID string `json:"accountId"`
		Nickname  string `json:"nickname"`
		Role      int    `json:"role"`
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
	m.AccountID = dt.Str(data["accountId"])
	m.Nickname = dt.Str(data["nickname"])
	m.Role = dt.Int(data["role"])
}
