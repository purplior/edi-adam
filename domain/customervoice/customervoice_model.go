package customervoice

import "time"

var (
	CustomerVoiceType_NewFeature CustomerVoiceType = "new_feature"
	CustomerVoiceType_BugReport  CustomerVoiceType = "bug_report"
	CustomerVoiceType_Withdrawal CustomerVoiceType = "withdrawal"
	CustomerVoiceType_Etc        CustomerVoiceType = "etc"
)

type (
	CustomerVoiceType string

	CustomerVoice struct {
		ID        string            `json:"id"`
		Type      CustomerVoiceType `json:"type"`
		UserID    string            `json:"userId"`
		Content   string            `json:"content"`
		CreatedAt time.Time         `json:"createdAt"`
	}
)

type (
	CustomerVoiceRegisterRequest struct {
		Type    CustomerVoiceType `json:"type"`
		Content string            `json:"content"`
	}
)
