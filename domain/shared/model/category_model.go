package model

import "time"

type (
	// @MODEL {카테고리}
	Category struct {
		ID        string    `gorm:"primaryKey;size:20" json:"id,omitempty"`
		Label     string    `gorm:"unique;size:50" json:"label"` // 한글 레이블
		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

		Assistants []Assistant `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"assistants,omitempty"`
	}

	// @MODEL {카테고리 칩}
	// - Assistant에서 사용되는 카테고리인지 정보 노출하지 않음
	// - ex. 카드 리스트
	CategoryChip struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	}
)

// @METHOD {칩으로 변환}
func (m Category) ToChip() CategoryChip {
	return CategoryChip{
		ID:    m.ID,
		Label: m.Label,
	}
}
