package challenge

import (
	"github.com/podossaem/podoroot/domain/ledger"
	"github.com/podossaem/podoroot/domain/mission"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/wallet"
	"github.com/podossaem/podoroot/lib/mydate"
)

type (
	ChallengeService interface {
		ReceiveOne(
			ctx inner.Context,
			id string,
			userID string,
		) error
	}
)

type (
	challengeService struct {
		challengeRepository ChallengeRepository
		walletService       wallet.WalletService
		cm                  inner.ContextManager
	}
)

func (s *challengeService) ReceiveOne(
	ctx inner.Context,
	id string,
	userID string,
) error {
	challenge, err := s.challengeRepository.FindOne_ByID(
		ctx,
		id,
	)

	if err != nil {
		if err == exception.ErrNoRecord {
			return exception.ErrBadRequest
		}
		return err
	}

	if challenge.UserID != userID {
		return exception.ErrBadRequest
	}

	if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	if err := s.challengeRepository.UpdateOne_ReceivedStatus_ByID(
		ctx,
		id,
		true,
		mydate.Now(),
	); err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	switch challenge.Mission.Reward {
	case mission.MissionReward_Podo5000:
		{
			if err := s.walletService.Charge(
				ctx,
				userID,
				5000,
				ledger.LedgerAction_ReceiveByMission,
				challenge.MissionID,
			); err != nil {
				s.cm.RollbackTX(ctx, inner.TX_PodoSql)
				return err
			}
		}
	}

	return s.cm.CommitTX(ctx, inner.TX_PodoSql)
}

func NewChallengeService(
	challengeRepository ChallengeRepository,
	walletService wallet.WalletService,
	cm inner.ContextManager,
) ChallengeService {
	return &challengeService{
		challengeRepository: challengeRepository,
		walletService:       walletService,
		cm:                  cm,
	}
}
