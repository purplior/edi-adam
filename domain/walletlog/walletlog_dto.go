package walletlog

import "github.com/purplior/edi-adam/domain/shared/model"

type (
	QueryOption struct {
		ID uint
	}
)

type (
	AddDTO struct {
		Type     model.WalletLogType `json:"type"`
		Delta    int                 `json:"delta"`
		Comment  string              `json:"comment"`
		WalletID uint                `json:"walletId"`
	}
)
