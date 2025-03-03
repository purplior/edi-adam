package customervoice

import (
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/dt"
)

type (
	QueryOption struct{}
)

type (
	CustomerVoiceRegisterDTO struct {
		UserID  uint
		Type    model.CustomerVoiceType `json:"type"`
		Content string                  `json:"content"`
	}
)

func (r CustomerVoiceRegisterDTO) ToSlackMessageText() string {
	text := "*[샘비서 문의]*"
	text += "\n*사용자 ID*: " + dt.Str(r.UserID)
	text += "\n*타입*: "
	switch r.Type {
	case model.CustomerVoiceType_NewFeature:
		text += "신기능 제안"
	case model.CustomerVoiceType_BugReport:
		text += "버그 제보"
	case model.CustomerVoiceType_Withdrawal:
		text += "회원 탈퇴"
	case model.CustomerVoiceType_Etc:
		text += "기타"
	}

	text += "\n*내용*:\n"
	text += r.Content

	return text
}
