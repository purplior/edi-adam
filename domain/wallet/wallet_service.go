package wallet

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	WalletService interface {
		RegisterOne(
			ctx context.APIContext,
			wallet Wallet,
		) (
			Wallet,
			error,
		)
	}
)

type (
	walletService struct {
		walletRepository WalletRepository
	}
)

func (s *walletService) RegisterOne(
	ctx context.APIContext,
	wallet Wallet,
) (
	Wallet,
	error,
) {
	return s.walletRepository.InsertOne(ctx, wallet)
}

func NewWalletService(
	walletRepository WalletRepository,
) WalletService {
	return &walletService{
		walletRepository: walletRepository,
	}
}
