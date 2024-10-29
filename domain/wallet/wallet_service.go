package wallet

import (
	"github.com/podossaem/podoroot/domain/ledger"
	"github.com/podossaem/podoroot/domain/shared/inner"
)

type (
	WalletService interface {
		RegisterOne(
			ctx inner.Context,
			wallet Wallet,
		) (
			Wallet,
			error,
		)

		Expend(
			ctx inner.Context,
			userId string,
			podoDelta int,
			ledgerAction ledger.LedgerAction,
			ledgerReason string,
		) error
	}
)

type (
	walletService struct {
		walletRepository WalletRepository
		ledgerService    ledger.LedgerService
	}
)

func (s *walletService) RegisterOne(
	ctx inner.Context,
	wallet Wallet,
) (
	Wallet,
	error,
) {
	return s.walletRepository.InsertOne(ctx, wallet)
}

func (s *walletService) Expend(
	ctx inner.Context,
	userId string,
	podoDelta int,
	ledgerAction ledger.LedgerAction,
	ledgerReason string,
) error {
	wallet, err := s.walletRepository.UpdateOneByUserIDAndDelta(
		ctx,
		userId,
		podoDelta,
	)
	if err != nil {
		return err
	}

	if _, err = s.ledgerService.RegisterOne(ctx, ledger.Ledger{
		WalletID:   wallet.ID,
		PodoAmount: podoDelta,
		Action:     ledgerAction,
		Reason:     ledgerReason,
	}); err != nil {
		return err
	}

	return nil
}

func NewWalletService(
	walletRepository WalletRepository,
	ledgerService ledger.LedgerService,
) WalletService {
	return &walletService{
		walletRepository: walletRepository,
		ledgerService:    ledgerService,
	}
}
