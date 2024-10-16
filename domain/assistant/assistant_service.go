package assistant

import "github.com/podossaem/podoroot/domain/context"

type (
	AssistantService interface {
		RegisterOne(
			ctx context.APIContext,
			authorID string,
			request RegisterOneRequest,
		) (
			Assistant,
			error,
		)

		GetList(
			ctx context.APIContext,
			authorID string,
		) (
			[]Assistant,
			error,
		)
	}
)

type (
	assistantService struct {
		assistantRepository AssistantRepository
	}
)

func (s *assistantService) RegisterOne(
	ctx context.APIContext,
	authorID string,
	request RegisterOneRequest,
) (
	Assistant,
	error,
) {
	assistant := Assistant{
		AuthorID:    authorID,
		Title:       request.Title,
		Description: request.Description,
		IsPublic:    request.IsPublic,
	}

	return s.assistantRepository.InsertOne(
		ctx,
		assistant,
	)
}

func (s *assistantService) GetList(
	ctx context.APIContext,
	authorID string,
) (
	[]Assistant,
	error,
) {
	return s.assistantRepository.FindListByAuthorID(
		ctx,
		authorID,
	)
}

func NewAssistantService(
	assistantRepository AssistantRepository,
) AssistantService {
	return &assistantService{
		assistantRepository: assistantRepository,
	}
}
