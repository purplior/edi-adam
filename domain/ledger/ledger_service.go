package ledger

import "github.com/purplior/podoroot/domain/shared/inner"

type (
	LedgerService interface {
		RegisterOne(
			ctx inner.Context,
			ledger Ledger,
		) (
			Ledger,
			error,
		)
	}
)

type (
	ledgerService struct {
		ledgerRepository LedgerRepository
	}
)

func (s *ledgerService) RegisterOne(
	ctx inner.Context,
	ledger Ledger,
) (
	Ledger,
	error,
) {
	return s.ledgerRepository.InsertOne(ctx, ledger)
}

func NewLedgerService(
	ledgerRepository LedgerRepository,
) LedgerService {
	return &ledgerService{
		ledgerRepository: ledgerRepository,
	}
}
