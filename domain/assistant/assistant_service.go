package assistant

import (
	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/assisterform"
	"github.com/purplior/podoroot/domain/shared/exception"
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

		UpdateOne(
			ctx inner.Context,
			authorID string,
			request UpdateOneRequest,
		) error

		RemoveOne_ByID(
			ctx inner.Context,
			authorID string,
			id string,
		) error

		GetOne_ByID(
			ctx inner.Context,
			id string,
			joinOption AssistantJoinOption,
		) (
			Assistant,
			error,
		)

		GetOne_ByViewID(
			ctx inner.Context,
			viewID string,
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

		ApproveOne(
			ctx inner.Context,
			id string,
			metaTags []string,
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

	status := AssistantStatus_Registered
	if request.IsPublic {
		status = AssistantStatus_UnderReview
	}

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
			IsPublic:      false,
			Status:        status,
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
			Version:            "v1.0.0",
			VersionDescription: "- 기본 기능 배포",
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

func (s *assistantService) UpdateOne(
	ctx inner.Context,
	authorID string,
	request UpdateOneRequest,
) error {
	if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
			panic(r)
		}
	}()

	_assistant, err := s.assistantRepository.FindOne_ByID(
		ctx,
		request.ID,
		AssistantJoinOption{
			WithAssisters: true,
		},
	)
	if err != nil {
		return err
	}
	if _assistant.IsPublic || len(_assistant.Assisters) == 0 {
		return exception.ErrBadRequest
	}
	if _assistant.AuthorID != authorID {
		return exception.ErrUnauthorized
	}

	_assister := _assistant.Assisters[0]
	_assisterForm, err := s.assisterformService.GetOne_ByAssisterID(
		ctx,
		_assister.ID,
	)
	if err != nil {
		return err
	}

	status := AssistantStatus_Registered
	if request.IsPublic {
		status = AssistantStatus_UnderReview
	}
	err = s.assistantRepository.UpdateOne(
		ctx,
		Assistant{
			ID:            request.ID,
			ViewID:        _assistant.ViewID,
			AuthorID:      authorID,
			CategoryID:    request.CategoryID,
			AssistantType: _assistant.AssistantType,
			Title:         request.Title,
			Description:   request.Description,
			Tags:          request.Tags,
			IsPublic:      false,
			Status:        status,
			CreatedAt:     _assistant.CreatedAt,
		},
	)
	if err != nil {
		return err
	}

	err = s.assisterformService.UpdateOne(
		ctx,
		assisterform.AssisterForm{
			ID:            _assisterForm.ID,
			AssisterID:    _assister.AssistantID,
			Origin:        _assisterForm.Origin,
			Model:         _assisterForm.Model,
			Fields:        request.Fields,
			Tests:         request.Tests,
			QueryMessages: request.QueryMessages,
			CreatedAt:     _assistant.CreatedAt,
		},
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	return nil
}

func (s *assistantService) RemoveOne_ByID(
	ctx inner.Context,
	authorID string,
	id string,
) error {
	assistant, err := s.assistantRepository.FindOne_ByID(
		ctx,
		id,
		AssistantJoinOption{
			WithAssisters: true,
		},
	)
	if err != nil {
		return err
	}
	if assistant.Status != AssistantStatus_Registered || assistant.AuthorID != authorID {
		return exception.ErrBadRequest
	}

	if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
			panic(r)
		}
	}()

	assisterIDs := make([]string, len(assistant.Assisters))
	for i, assister := range assistant.Assisters {
		assisterIDs[i] = assister.ID
	}

	if err := s.assisterService.RemoveAll_ByIDs(
		ctx,
		assisterIDs,
	); err != nil {
		return err
	}

	if err := s.assisterformService.RemoveAll_ByAssisterIDs(
		ctx,
		assisterIDs,
	); err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	if err := s.assistantRepository.DeleteOne_ByID(
		ctx,
		id,
	); err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	return nil
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

func (s *assistantService) GetOne_ByViewID(
	ctx inner.Context,
	viewID string,
	joinOption AssistantJoinOption,
) (
	Assistant,
	error,
) {
	return s.assistantRepository.FindOne_ByViewID(ctx, viewID, joinOption)
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
		assistantInfos[i] = assistant.ToInfo()
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

func (s *assistantService) ApproveOne(
	ctx inner.Context,
	id string,
	metaTags []string,
) error {
	assistant, err := s.assistantRepository.FindOne_ByID(
		ctx,
		id,
		AssistantJoinOption{},
	)
	if err != nil {
		return err
	}
	if assistant.Status != AssistantStatus_UnderReview {
		return exception.ErrBadRequest
	}

	assistant.Status = AssistantStatus_Approved
	assistant.MetaTags = metaTags
	assistant.IsPublic = true
	err = s.assistantRepository.UpdateOne(
		ctx,
		assistant,
	)
	if err != nil {
		return err
	}

	return nil
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
