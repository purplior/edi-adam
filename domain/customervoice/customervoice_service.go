package customervoice

import "github.com/purplior/podoroot/domain/shared/inner"

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
	return s.customerVoiceRepository.InsertOne(ctx, CustomerVoice{
		UserID:  userID,
		Title:   request.Title,
		Content: request.Content,
	})
}

func NewCustomerVoiceService(
	customerVoiceRepository CustomerVoiceRepository,
) CustomerVoiceService {
	return &customerVoiceService{
		customerVoiceRepository: customerVoiceRepository,
	}
}
