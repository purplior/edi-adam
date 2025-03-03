package missionlog

import (
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/domain/wallet"
	"github.com/purplior/edi-adam/domain/walletlog"
	"github.com/purplior/edi-adam/lib/mydate"
)

type (
	MissionLogService interface {
		GetList(
			session inner.Session,
			queryOption QueryOption,
		) (
			[]model.MissionLog,
			error,
		)

		Receive(
			session inner.Session,
			id uint,
			userId uint,
		) error

		Achieve(
			session inner.Session,
			dto AchieveDTO,
		) error
	}
)

type (
	missionLogService struct {
		missionLogRepository MissionLogRepository
		walletService        wallet.WalletService
	}
)

func (s *missionLogService) GetList(
	session inner.Session,
	queryOption QueryOption,
) (
	[]model.MissionLog,
	error,
) {
	return s.missionLogRepository.ReadList(
		session,
		queryOption,
	)
}

func (s *missionLogService) Receive(
	session inner.Session,
	id uint,
	userId uint,
) (err error) {
	var m model.MissionLog
	m, err = s.missionLogRepository.Read(
		session,
		QueryOption{
			ID:          id,
			UserID:      session.Identity().ID,
			WithMission: true,
		},
	)
	if err != nil || m.Mission == nil {
		return err
	}
	if m.UserID != userId || !m.IsAchieved {
		return exception.ErrBadRequest
	}
	if m.IsReceived {
		return exception.ErrAlreadyReceived
	}

	if err = session.BeginTransaction(); err != nil {
		return err
	}

	receivedAt := mydate.Now()
	if err = s.missionLogRepository.Updates(
		session,
		QueryOption{
			ID: m.ID,
		},
		model.MissionLog{
			IsReceived: true,
			ReceivedAt: &receivedAt,
		},
	); err != nil {
		return err
	}

	switch m.Mission.RewardType {
	case model.MissionRewardType_Coin:
		err = s.walletService.Charge(session, wallet.ChargeDTO{
			OwnerID:   userId,
			Delta:     uint(m.Mission.RewardAmount),
			LogAddDTO: walletlog.AddDTO{},
		})
	}
	if err != nil {
		return err
	}
	if err = session.CommitTransaction(); err != nil {
		return err
	}

	return nil
}

func (s *missionLogService) Achieve(
	session inner.Session,
	dto AchieveDTO,
) (err error) {
	var m model.MissionLog
	m, err = s.missionLogRepository.Read(
		session,
		QueryOption{
			UserID:    dto.UserID,
			MissionID: dto.MissionID,
		},
	)
	isNoRecord := err == exception.ErrNoRecord
	if err != nil && !isNoRecord {
		return err
	}

	if isNoRecord {
		m, err = s.missionLogRepository.Create(
			session,
			model.MissionLog{
				IsAchieved: true,
				IsReceived: false,
				UserID:     dto.UserID,
				MissionID:  dto.MissionID,
			},
		)

		return err
	}

	// 이미 달성했거나 수령한 경우 무시함.
	if m.IsAchieved || m.IsReceived {
		return exception.ErrBadRequest
	}

	return s.missionLogRepository.Updates(
		session,
		QueryOption{
			ID: m.ID,
		},
		model.MissionLog{
			IsAchieved: true,
		},
	)
}

func NewMissionLogService(
	missionLogRepository MissionLogRepository,
) MissionLogService {
	return &missionLogService{
		missionLogRepository: missionLogRepository,
	}
}
