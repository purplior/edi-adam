package wallet

import "github.com/purplior/edi-adam/domain/walletlog"

type (
	QueryOption struct {
		OwnerID uint
	}

	AddDTO struct {
		OwnerID uint
	}

	ExpendDTO struct {
		OwnerID   uint             `json:"ownerId"`
		Delta     uint             `json:"delta"`
		LogAddDTO walletlog.AddDTO `json:"logAddDto"`
	}

	ChargeDTO struct {
		OwnerID   uint             `json:"ownerId"`
		Delta     uint             `json:"delta"`
		LogAddDTO walletlog.AddDTO `json:"logAddDto"`
	}
)
