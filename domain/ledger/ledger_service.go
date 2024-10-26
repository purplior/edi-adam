package ledger

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	LedgerService interface {
		RegisterOne(
			ctx context.APIContext,
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
	ctx context.APIContext,
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
