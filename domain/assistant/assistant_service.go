package assistant

import (
	"github.com/podossaem/podoroot/domain/shared/inner"
)

type (
	AssistantService interface {
		RegisterOne(
			ctx inner.Context,
			authorID string,
			request RegisterOneRequest,
		) (
			Assistant,
			error,
		)

		GetDetailOneByViewID(
			ctx inner.Context,
			viewID string,
			joinOption AssistantJoinOption,
		) (
			AssistantDetail,
			error,
		)

		GetInfoListByAuthor(
			ctx inner.Context,
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
	ctx inner.Context,
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

func (s *assistantService) GetDetailOneByViewID(
	ctx inner.Context,
	viewID string,
	joinOption AssistantJoinOption,
) (
	AssistantDetail,
	error,
) {
	assistant, err := s.assistantRepository.FindOneByViewID(
		ctx,
		viewID,
		joinOption,
	)
	if err != nil {
		return AssistantDetail{}, err
	}

	return assistant.ToDetail()
}

func (s *assistantService) GetInfoListByAuthor(
	ctx inner.Context,
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
