package model

import "time"

type (
	Bookmark struct {
		ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`

		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

		UserID      uint       `gorm:"uniqueIndex:idx_user_assistant" json:"userId,omitempty"`
		AssistantID uint       `gorm:"uniqueIndex:idx_user_assistant" json:"assistantId,omitempty"`
		Assistant   *Assistant `json:"assistant,omitempty"`
	}
)
