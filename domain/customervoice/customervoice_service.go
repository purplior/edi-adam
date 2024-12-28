package customervoice

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/user"
)

type (
	CustomerVoiceService interface {
		RegisterOne(
			ctx inner.Context,
			userID string,
			request CustomerVoiceRegisterRequest,
		) (
			CustomerVoice,
			error,
		)
	}
)

type (
	customerVoiceService struct {
		customerVoiceRepository CustomerVoiceRepository
		userService             user.UserService
		cm                      inner.ContextManager
	}
)

func (s *customerVoiceService) RegisterOne(
	ctx inner.Context,
	userID string,
	request CustomerVoiceRegisterRequest,
) (
	CustomerVoice,
	error,
) {
	if request.Type == CustomerVoiceType_Withdrawal {
		return s.processWithdrawal(
			ctx,
			userID,
			request,
		)
	}

	return s.customerVoiceRepository.InsertOne(ctx, CustomerVoice{
		UserID:  userID,
		Type:    request.Type,
		Content: request.Content,
	})
}

func (s *customerVoiceService) processWithdrawal(
	ctx inner.Context,
	userID string,
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
		UserID:  userID,
		Type:    request.Type,
		Content: request.Content,
	})
	if err != nil {
		return CustomerVoice{}, err
	}
	if err := s.userService.Inactive(ctx, userID); err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return CustomerVoice{}, err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return CustomerVoice{}, err
	}

	return customerVoice, nil
}

func NewCustomerVoiceService(
	customerVoiceRepository CustomerVoiceRepository,
	userService user.UserService,
	cm inner.ContextManager,
) CustomerVoiceService {
	return &customerVoiceService{
		customerVoiceRepository: customerVoiceRepository,
		userService:             userService,
		cm:                      cm,
	}
}
