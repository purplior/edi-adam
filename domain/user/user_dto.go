package user

import (
	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/security"
)

type (
	QueryOption struct {
		ID          uint
		AccountID   string
		PhoneNumber string
		Nickname    string
	}
)

func PhoneNumberToAccountID(phoneNumber string) (string, error) {
	if len(phoneNumber) == 0 {
		return "", nil
	}

	accountID, err := security.EncryptMapDataWithAESGCM(
		map[string]interface{}{
			"method":       model.UserJoinMethod_Phone,
			"phone_number": phoneNumber,
		},
		config.SymmetricKey(),
	)
	if err != nil {
		return "", err
	}

	return accountID, nil
}

type (
	// @MODEL {사용자 등록요청 모델}
	RegisterMemberDTO struct {
		PhoneNumber      string           `json:"phoneNumber"`
		Nickname         string           `json:"nickname"`
		Avatar           model.UserAvatar `json:"avatar"`
		IsMarketingAgree bool             `json:"isMarketingAgree"`
	}
)
