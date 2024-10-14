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
		AuthorID:     authorID,
		Title:        request.Title,
		Description:  request.Description,
		VersionLabel: request.VersionLabel,
		IsPublic:     request.IsPublic,
	}

	return s.assistantRepository.InsertOne(
		ctx,
		assistant,
	)
}

func NewAssistantService(
	assistantRepository AssistantRepository,
) AssistantService {
	return &assistantService{
		assistantRepository: assistantRepository,
	}
}
