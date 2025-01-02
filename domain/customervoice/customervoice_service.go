package customervoice

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/logger"
	"github.com/purplior/podoroot/domain/user"
	"github.com/purplior/podoroot/infra/port/podoslack"
)

type (
	CustomerVoiceService interface {
		RegisterOne(
			ctx inner.Context,
			request CustomerVoiceRegisterRequest,
		) (
			CustomerVoice,
			error,
		)
	}
)

type (
	customerVoiceService struct {
		slackClient             *podoslack.Client
		customerVoiceRepository CustomerVoiceRepository
		userService             user.UserService
		cm                      inner.ContextManager
	}
)

func (s *customerVoiceService) RegisterOne(
	ctx inner.Context,
	request CustomerVoiceRegisterRequest,
) (
	CustomerVoice,
	error,
) {
	if err := s.slackClient.SendMessage(podoslack.SendMessageRequest{
		ChannelID: podoslack.ChannelID_CustomerVoice,
		Text:      request.ToSlackMessageText(),
	}); err != nil {
		logger.Error(err, "slack error occurred")
	}

	if request.Type == CustomerVoiceType_Withdrawal {
		return s.processWithdrawal(
			ctx,
			request,
		)
	}

	return s.customerVoiceRepository.InsertOne(ctx, CustomerVoice{
		UserID:  request.UserID,
		Type:    request.Type,
		Content: request.Content,
	})
}

func (s *customerVoiceService) processWithdrawal(
	ctx inner.Context,
	request CustomerVoiceRegisterRequest,
) (
	CustomerVoice,
	error,
) {
	if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
		return CustomerVoice{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
			panic(r)
		}
	}()

	customerVoice, err := s.customerVoiceRepository.InsertOne(ctx, CustomerVoice{
		UserID:  request.UserID,
		Type:    request.Type,
		Content: request.Content,
	})
	if err != nil {
		return CustomerVoice{}, err
	}
	if err := s.userService.Inactive(ctx, request.UserID); err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return CustomerVoice{}, err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return CustomerVoice{}, err
	}

	return customerVoice, nil
}

func NewCustomerVoiceService(
	slackClient *podoslack.Client,
	customerVoiceRepository CustomerVoiceRepository,
	userService user.UserService,
	cm inner.ContextManager,
) CustomerVoiceService {
	return &customerVoiceService{
		slackClient:             slackClient,
		customerVoiceRepository: customerVoiceRepository,
		userService:             userService,
		cm:                      cm,
	}
}
