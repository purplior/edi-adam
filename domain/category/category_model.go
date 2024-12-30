package category

import "time"

type (
	Category struct {
		ID        string    `json:"id"`
		Alias     string    `json:"alias"`
		Label     string    `json:"label"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

func (m Category) ToInfo() CategoryInfo {
	return CategoryInfo{
		ID:    m.ID,
		Alias: m.Alias,
		Label: m.Label,
	}
}

type (
	CategoryInfo struct {
		ID    string `json:"id"`
		Alias string `json:"alias"`
		Label string `json:"label"`
	}
)
