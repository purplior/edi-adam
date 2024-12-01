package category

import "time"

type (
	Category struct {
		ID        string    `json:"id"`
		Alias     string    `json:"alias"`
		Label     string    `json:"label"`
		CreatorID string    `json:"creatorId"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

func (m Category) ToInfo() CategoryInfo {
	return CategoryInfo{
		Alias: m.Alias,
		Label: m.Label,
	}
}

type (
	CategoryInfo struct {
		Alias string `json:"alias"`
		Label string `json:"label"`
	}
)
