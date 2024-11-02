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

		GetOne_ByUserID(
			ctx inner.Context,
			userID string,
		) (
			Wallet,
			error,
		)

		Expend(
			ctx inner.Context,
			userID string,
			podoDelta int,
			ledgerAction ledger.LedgerAction,
			ledgerReason string,
		) error

		Charge(
			ctx inner.Context,
			userID string,
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

func (s *walletService) GetOne_ByUserID(
	ctx inner.Context,
	userID string,
) (
	Wallet,
	error,
) {
	return s.walletRepository.FindOne_ByUserID(
		ctx,
		userID,
	)
}

func (s *walletService) Expend(
	ctx inner.Context,
	userID string,
	podoDelta int,
	ledgerAction ledger.LedgerAction,
	ledgerReason string,
) error {
	// 0은 히스토리에 남기지 않아요
	if podoDelta == 0 {
		return nil
	}

	podoDelta = -1 * podoDelta
	wallet, err := s.walletRepository.UpdateOne_ByUserIDAndDelta(
		ctx,
		userID,
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

func (s *walletService) Charge(
	ctx inner.Context,
	userID string,
	podoDelta int,
	ledgerAction ledger.LedgerAction,
	ledgerReason string,
) error {
	// 0은 히스토리에 남기지 않아요
	if podoDelta == 0 {
		return nil
	}

	wallet, err := s.walletRepository.UpdateOne_ByUserIDAndDelta(
		ctx,
		userID,
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
