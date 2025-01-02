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
		UserID  string
		Type    CustomerVoiceType `json:"type"`
		Content string            `json:"content"`
	}
)

func (r CustomerVoiceRegisterRequest) ToSlackMessageText() string {
	text := "*[포도쌤 문의]*"
	text += "\n*사용자 ID*: " + r.UserID
	text += "\n*타입*: "
	switch r.Type {
	case CustomerVoiceType_NewFeature:
		text += "신기능 제안"
	case CustomerVoiceType_BugReport:
		text += "버그 제보"
	case CustomerVoiceType_Withdrawal:
		text += "회원 탈퇴"
	case CustomerVoiceType_Etc:
		text += "기타"
	}

	text += "\n*내용*:\n"
	text += r.Content

	return text
}
