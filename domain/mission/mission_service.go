package mission

import (
	"strings"

	"github.com/purplior/sbec/domain/challenge"
	"github.com/purplior/sbec/domain/ledger"
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/shared/pagination"
	"github.com/purplior/sbec/domain/wallet"
	"github.com/purplior/sbec/lib/dt"
)

type (
	MissionService interface {
		GetPaginatedList_OnlyPublic_ByUserID(
			ctx inner.Context,
			userID string,
			page int,
			pageSize int,
		) (
			[]Mission,
			pagination.PaginationMeta,
			error,
		)

		ReceiveOne(
			ctx inner.Context,
			id string,
			userID string,
		) error
	}
)

type (
	missionService struct {
		missionRepository MissionRepository
		challengeService  challenge.ChallengeService
		walletService     wallet.WalletService
		cm                inner.ContextManager
	}
)

func (s *missionService) GetPaginatedList_OnlyPublic_ByUserID(
	ctx inner.Context,
	userID string,
	page int,
	pageSize int,
) (
	[]Mission,
	pagination.PaginationMeta,
	error,
) {
	missions, pageMeta, err := s.missionRepository.FindPaginatedList_OnlyPublic_ByUserID(
		ctx,
		userID,
		page,
		pageSize,
	)
	if err != nil {
		return nil, pagination.PaginationMeta{}, err
	}

	missionIDs := make([]string, len(missions))
	for i, mission := range missions {
		missionIDs[i] = mission.ID
	}

	challenges, err := s.challengeService.GetList_ByUserIDAndMissionIDs(
		ctx,
		userID,
		missionIDs,
	)
	if err != nil {
		return nil, pagination.PaginationMeta{}, err
	}

	if len(challenges) > 0 {
		j := 0
		for i, mission := range missions {
			if mission.ID != challenges[j].MissionID {
				continue
			}

			missions[i].Challenges = []challenge.Challenge{
				challenges[j],
			}
			j++

			if j >= len(challenges) {
				break
			}
		}
	}

	return missions, pageMeta, err
}

func (s *missionService) ReceiveOne(
	ctx inner.Context,
	id string,
	userID string,
) error {
	mission, err := s.missionRepository.FindOne_ByIDAndUserID(
		ctx,
		id,
		userID,
	)

	if err != nil {
		if err == exception.ErrNoRecord {
			return exception.ErrBadRequest
		}
		return err
	}
	if len(mission.Challenges) == 0 {
		return exception.ErrBadRequest
	}
	if !mission.Challenges[0].IsAchieved {
		return exception.ErrBadRequest
	}
	if mission.Challenges[0].IsReceived {
		return exception.ErrAlreadyReceived
	}

	if err := s.cm.BeginTX(ctx, inner.TX_sqldb); err != nil {
		return err
	}

	if err := s.challengeService.PatchOne_ReceivedStatus(
		ctx,
		mission.Challenges[0].ID,
	); err != nil {
		s.cm.RollbackTX(ctx, inner.TX_sqldb)
		return err
	}

	reward := string(mission.Reward)

	if strings.HasPrefix(reward, "podo_") {
		podoAmount := dt.Int(strings.ReplaceAll(reward, "podo_", ""))

		if err := s.walletService.Charge(
			ctx,
			userID,
			podoAmount,
			ledger.LedgerAction_ReceiveByMission,
			mission.ID,
		); err != nil {
			s.cm.RollbackTX(ctx, inner.TX_sqldb)
			return err
		}
	}

	return s.cm.CommitTX(ctx, inner.TX_sqldb)
}

func NewMissionService(
	missionRepository MissionRepository,
	challengeService challenge.ChallengeService,
	walletService wallet.WalletService,
	cm inner.ContextManager,
) MissionService {
	return &missionService{
		missionRepository: missionRepository,
		challengeService:  challengeService,
		walletService:     walletService,
		cm:                cm,
	}
}
