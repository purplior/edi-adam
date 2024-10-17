package assistant

import "github.com/podossaem/podoroot/domain/shared/context"

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

		GetDetailOneByID(
			ctx context.APIContext,
			id string,
			joinOption AssistantJoinOption,
		) (
			AssistantDetail,
			error,
		)

		GetInfoListByAuthor(
			ctx context.APIContext,
			authorID string,
			joinOption AssistantJoinOption,
		) (
			[]AssistantInfo,
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

func (s *assistantService) GetDetailOneByID(
	ctx context.APIContext,
	id string,
	joinOption AssistantJoinOption,
) (
	AssistantDetail,
	error,
) {
	assistant, err := s.assistantRepository.FindOneByID(
		ctx,
		id,
		joinOption,
	)
	if err != nil {
		return AssistantDetail{}, err
	}

	return assistant.ToDetail()
}

func (s *assistantService) GetInfoListByAuthor(
	ctx context.APIContext,
	authorID string,
	joinOption AssistantJoinOption,
) (
	[]AssistantInfo,
	error,
) {
	assistants, err := s.assistantRepository.FindListByAuthorID(
		ctx,
		authorID,
		joinOption,
	)
	if err != nil {
		return nil, err
	}

	assistantInfos := make([]AssistantInfo, len(assistants))
	for i, assistant := range assistants {
		assistantInfos[i], err = assistant.ToInfo()
		if err != nil {
			return nil, err
		}
	}

	return assistantInfos, nil
}

func NewAssistantService(
	assistantRepository AssistantRepository,
) AssistantService {
	return &assistantService{
		assistantRepository: assistantRepository,
	}
}
