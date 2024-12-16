package assistant

import (
	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/assisterform"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/lib/strgen"
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

		GetOne_ByID(
			ctx inner.Context,
			id string,
			joinOption AssistantJoinOption,
		) (
			Assistant,
			error,
		)

		GetDetailOne_ByViewID(
			ctx inner.Context,
			viewID string,
			joinOption AssistantJoinOption,
		) (
			AssistantDetail,
			error,
		)

		GetInfoList_ByAuthor(
			ctx inner.Context,
			authorID string,
			joinOption AssistantJoinOption,
		) (
			[]AssistantInfo,
			error,
		)

		GetInfoList_ByCategory(
			ctx inner.Context,
			categoryAlias string,
			joinOption AssistantJoinOption,
		) (
			[]AssistantInfo,
			error,
		)

		GetPaginatedList_ByAuthor(
			ctx inner.Context,
			authorID string,
			page int,
			pageSize int,
		) (
			[]Assistant,
			pagination.PaginationMeta,
			error,
		)

		PutOne(
			ctx inner.Context,
			assistant Assistant,
		) error

		CreateOne(
			ctx inner.Context,
			assistant Assistant,
		) error
	}
)

type (
	assistantService struct {
		assistantRepository AssistantRepository
		assisterService     assister.AssisterService
		assisterformService assisterform.AssisterFormService
		cm                  inner.ContextManager
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
	if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
		return Assistant{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
			panic(r)
		}
	}()

	_assistant, err := s.assistantRepository.InsertOne(
		ctx,
		Assistant{
			ViewID:        strgen.ShortUniqueID(),
			AuthorID:      authorID,
			CategoryID:    request.CategoryID,
			AssistantType: AssistantType_Formal,
			Title:         request.Title,
			Description:   request.Description,
			Tags:          request.Tags,
			IsPublic:      request.IsPublic,
		},
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return Assistant{}, err
	}

	_assister, err := s.assisterService.RegisterOne(
		ctx,
		assister.AssisterRegisterRequest{
			AssistantID:        _assistant.ID,
			Version:            request.Version,
			VersionDescription: request.VersionDescription,
			Cost:               2,
		},
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return Assistant{}, err
	}

	_, err = s.assisterformService.RegisterOne(
		ctx,
		assisterform.AssisterFormRegisterRequest{
			AssisterID:    _assister.ID,
			Origin:        assisterform.AssisterOrigin_OpenAI,
			Model:         assisterform.AssisterModel_OpenAI_ChatGPT4o,
			Tests:         request.Tests,
			Fields:        request.Fields,
			QueryMessages: request.QueryMessages,
		},
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return Assistant{}, err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return Assistant{}, err
	}

	_assistant.Assisters = []assister.Assister{
		_assister,
	}

	return _assistant, nil
}

func (s *assistantService) GetOne_ByID(
	ctx inner.Context,
	id string,
	joinOption AssistantJoinOption,
) (
	Assistant,
	error,
) {
	return s.assistantRepository.FindOne_ByID(ctx, id, joinOption)
}

func (s *assistantService) GetDetailOne_ByViewID(
	ctx inner.Context,
	viewID string,
	joinOption AssistantJoinOption,
) (
	AssistantDetail,
	error,
) {
	assistant, err := s.assistantRepository.FindOne_ByViewID(
		ctx,
		viewID,
		joinOption,
	)
	if err != nil {
		return AssistantDetail{}, err
	}

	return assistant.ToDetail()
}

func (s *assistantService) GetInfoList_ByAuthor(
	ctx inner.Context,
	authorID string,
	joinOption AssistantJoinOption,
) (
	[]AssistantInfo,
	error,
) {
	assistants, err := s.assistantRepository.FindList_ByAuthorID(
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

func (s *assistantService) GetInfoList_ByCategory(
	ctx inner.Context,
	categoryAlias string,
	joinOption AssistantJoinOption,
) (
	[]AssistantInfo,
	error,
) {
	assistants, err := s.assistantRepository.FindList_ByCategoryAlias(
		ctx,
		categoryAlias,
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

func (s *assistantService) GetPaginatedList_ByAuthor(
	ctx inner.Context,
	authorID string,
	page int,
	pageSize int,
) (
	[]Assistant,
	pagination.PaginationMeta,
	error,
) {
	return s.assistantRepository.FindPaginatedList_ByAuthorID(
		ctx,
		authorID,
		page,
		pageSize,
	)
}

func (s *assistantService) PutOne(
	ctx inner.Context,
	assistant Assistant,
) error {
	return s.assistantRepository.UpdateOne(ctx, assistant)
}

func (s *assistantService) CreateOne(
	ctx inner.Context,
	assistant Assistant,
) error {
	assistant.ViewID = strgen.ShortUniqueID()

	_, err := s.assistantRepository.InsertOne(
		ctx,
		assistant,
	)

	return err
}

func NewAssistantService(
	assistantRepository AssistantRepository,
	assisterService assister.AssisterService,
	assisterformService assisterform.AssisterFormService,
	cm inner.ContextManager,
) AssistantService {
	return &assistantService{
		assistantRepository: assistantRepository,
		assisterService:     assisterService,
		assisterformService: assisterformService,
		cm:                  cm,
	}
}
