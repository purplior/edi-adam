package wallet

import (
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/domain/walletlog"
)

type (
	WalletService interface {
		Get(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Wallet,
			error,
		)

		Add(
			session inner.Session,
			dto AddDTO,
		) (
			model.Wallet,
			error,
		)

		Expend(
			session inner.Session,
			dto ExpendDTO,
		) error

		Charge(
			session inner.Session,
			dto ChargeDTO,
		) error
	}
)

type (
	walletService struct {
		walletRepository WalletRepository
		walletLogService walletlog.WalletLogService
	}
)

func (s *walletService) Get(
	session inner.Session,
	queryOption QueryOption,
) (
	m model.Wallet,
	err error,
) {
	return s.walletRepository.Read(
		session,
		queryOption,
	)
}

func (s *walletService) Add(
	session inner.Session,
	dto AddDTO,
) (
	model.Wallet,
	error,
) {
	return s.walletRepository.Create(session, model.Wallet{
		OwnerID: dto.OwnerID,
		Coin:    0,
	})
}

// dto의 WalletID는 무시됨
func (s *walletService) Expend(
	session inner.Session,
	dto ExpendDTO,
) error {
	if dto.Delta == 0 {
		return nil
	}

	m, err := s.walletRepository.UpdatesCoinDelta_ByOwnerID(
		session,
		dto.OwnerID,
		-1*int(dto.Delta),
	)
	if err != nil {
		return err
	}

	if _, err = s.walletLogService.Add(session, walletlog.AddDTO{
		Type:     dto.LogAddDTO.Type,
		Delta:    int(dto.Delta),
		Comment:  dto.LogAddDTO.Comment,
		WalletID: m.ID,
	}); err != nil {
		return err
	}

	return nil
}

// dto의 WalletID는 무시됨
func (s *walletService) Charge(
	session inner.Session,
	dto ChargeDTO,
) error {
	if dto.Delta == 0 {
		return nil
	}

	m, err := s.walletRepository.UpdatesCoinDelta_ByOwnerID(
		session,
		dto.OwnerID,
		int(dto.Delta),
	)
	if err != nil {
		return err
	}

	if _, err = s.walletLogService.Add(session, walletlog.AddDTO{
		Type:     dto.LogAddDTO.Type,
		Delta:    int(dto.Delta),
		Comment:  dto.LogAddDTO.Comment,
		WalletID: m.ID,
	}); err != nil {
		return err
	}

	return nil
}

func NewWalletService(
	walletRepository WalletRepository,
	walletLogService walletlog.WalletLogService,
) WalletService {
	return &walletService{
		walletRepository: walletRepository,
		walletLogService: walletLogService,
	}
}
