package walletlog

import (
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	WalletLogService interface {
		Add(
			session inner.Session,
			dto AddDTO,
		) (
			model.WalletLog,
			error,
		)
	}
)

type (
	walletLogService struct {
		walletLogRepository WalletLogRepository
	}
)

func (s *walletLogService) Add(
	session inner.Session,
	dto AddDTO,
) (
	m model.WalletLog,
	err error,
) {
	return s.walletLogRepository.Create(
		session,
		model.WalletLog{
			Type:     dto.Type,
			Delta:    dto.Delta,
			Comment:  dto.Comment,
			WalletID: dto.WalletID,
		},
	)
}

func NewWalletLogService(
	walletLogRepository WalletLogRepository,
) WalletLogService {
	return &walletLogService{
		walletLogRepository: walletLogRepository,
	}
}
