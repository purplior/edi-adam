package model

import "time"

type (
	CustomerVoiceType string
)

var (
	CustomerVoiceType_NewFeature CustomerVoiceType = "new_feature"
	CustomerVoiceType_BugReport  CustomerVoiceType = "bug_report"
	CustomerVoiceType_Withdrawal CustomerVoiceType = "withdrawal"
	CustomerVoiceType_Etc        CustomerVoiceType = "etc"
)

type (
	CustomerVoice struct {
		ID      uint              `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		Type    CustomerVoiceType `gorm:"size:80" json:"type"`
		Content string            `gorm:"size:2000" json:"content"`

		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

		UserID uint `json:"userId,omitempty"`
	}
)
