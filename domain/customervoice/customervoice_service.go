package customervoice

import (
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/logger"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/domain/user"
	"github.com/purplior/edi-adam/infra/port/slack"
)

type (
	CustomerVoiceService interface {
		Register(
			session inner.Session,
			request CustomerVoiceRegisterDTO,
		) (
			model.CustomerVoice,
			error,
		)
	}
)

type (
	customerVoiceService struct {
		slackClient             *slack.Client
		customerVoiceRepository CustomerVoiceRepository
		userService             user.UserService
	}
)

func (s *customerVoiceService) Register(
	session inner.Session,
	dto CustomerVoiceRegisterDTO,
) (
	model.CustomerVoice,
	error,
) {
	if err := s.slackClient.SendMessage(slack.SendMessageRequest{
		ChannelID: slack.ChannelID_CustomerVoice,
		Text:      dto.ToSlackMessageText(),
	}); err != nil {
		logger.Error(err, "slack error occurred")
	}

	if dto.Type == model.CustomerVoiceType_Withdrawal {
		return s.processWithdrawal(
			session,
			dto,
		)
	}

	return s.create(session, dto)
}

func (s *customerVoiceService) processWithdrawal(
	session inner.Session,
	dto CustomerVoiceRegisterDTO,
) (
	m model.CustomerVoice,
	err error,
) {
	if err = session.BeginTransaction(); err != nil {
		return m, err
	}

	m, err = s.create(session, dto)
	if err != nil {
		return m, err
	}
	if err = s.userService.Inactive(session, dto.UserID); err != nil {
		session.RollbackTransaction()
		return m, err
	}
	if err = session.CommitTransaction(); err != nil {
		return m, err
	}

	return m, nil
}

func (s *customerVoiceService) create(
	session inner.Session,
	dto CustomerVoiceRegisterDTO,
) (model.CustomerVoice, error) {
	return s.customerVoiceRepository.Create(session, model.CustomerVoice{
		UserID:  dto.UserID,
		Type:    dto.Type,
		Content: dto.Content,
	})
}

func NewCustomerVoiceService(
	slackClient *slack.Client,
	customerVoiceRepository CustomerVoiceRepository,
	userService user.UserService,
) CustomerVoiceService {
	return &customerVoiceService{
		slackClient:             slackClient,
		customerVoiceRepository: customerVoiceRepository,
		userService:             userService,
	}
}
