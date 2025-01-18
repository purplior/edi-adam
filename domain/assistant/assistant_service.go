package assistant

import (
	"github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/domain/review"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/domain/wallet"
	"github.com/purplior/podoroot/infra/port/podoopenai"
	"github.com/purplior/podoroot/lib/dt"
	"github.com/purplior/podoroot/lib/mydate"
)

type (
	AssistantService interface {
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
		) (
			AssistantDetail,
			error,
		)

		GetPaginatedList_ByCategoryAlias(
			ctx inner.Context,
			categoryAlias string,
			isPublicOnly bool,
			pageRequest pagination.PaginationRequest,
		) (
			[]Assistant,
			pagination.PaginationMeta,
			error,
		)

		GetPaginatedList_ByAuthor(
			ctx inner.Context,
			authorID string,
			pageRequest pagination.PaginationRequest,
		) (
			[]Assistant,
			pagination.PaginationMeta,
			error,
		)

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

		ApproveOne(
			ctx inner.Context,
			id string,
			cost int,
			metaTags []string,
		) error
	}
)

type (
	assistantService struct {
		openaiClient        *podoopenai.Client
		assistantRepository AssistantRepository
		categoryService     category.CategoryService
		assisterService     assister.AssisterService
		walletService       wallet.WalletService
		reviewService       review.ReviewService
		cm                  inner.ContextManager
	}
)

func (s *assistantService) GetOne_ByID(
	ctx inner.Context,
	id string,
	joinOption AssistantJoinOption,
) (
	Assistant,
	error,
) {
	_assistant, err := s.assistantRepository.FindOne_ByID(ctx, id, joinOption)
	if err != nil {
		return Assistant{}, err
	}

	if joinOption.WithAssister {
		_assister, err := s.assisterService.GetOne_ByID(
			ctx,
			_assistant.AssisterID,
		)
		if err != nil {
			return Assistant{}, err
		}

		_assistant.Assister = _assister
	}

	return _assistant, nil
}

func (s *assistantService) GetOne_ByViewID(
	ctx inner.Context,
	viewID string,
	joinOption AssistantJoinOption,
) (
	Assistant,
	error,
) {
	_assistant, err := s.assistantRepository.FindOne_ByViewID(ctx, viewID, joinOption)
	if err != nil {
		return Assistant{}, err
	}

	if joinOption.WithAssister {
		_assister, err := s.assisterService.GetOne_ByID(
			ctx,
			_assistant.AssisterID,
		)
		if err != nil {
			return Assistant{}, err
		}

		_assistant.Assister = _assister
	}

	return _assistant, nil
}

func (s *assistantService) GetDetailOne_ByViewID(
	ctx inner.Context,
	viewID string,
) (
	AssistantDetail,
	error,
) {
	_assistant, err := s.assistantRepository.FindOne_ByViewID(ctx, viewID, AssistantJoinOption{
		WithAuthor:   true,
		WithAssister: true,
	})
	if err != nil {
		return AssistantDetail{}, err
	}

	_assister, err := s.assisterService.GetOne_ByID(
		ctx,
		_assistant.AssisterID,
	)
	if err != nil {
		return AssistantDetail{}, err
	}
	_assistant.Assister = _assister

	reviews, reviewPageMeta, err := s.reviewService.GetPaginatedList_ByAssistantID(
		ctx,
		_assistant.ID,
		pagination.PaginationRequest{
			Page: 1,
			Size: 10,
		},
	)
	if err != nil && err != exception.ErrNoRecord {
		return AssistantDetail{}, err
	}
	_assistant.Reviews = reviews
	_assistant.ReviewPageMeta = reviewPageMeta

	assistantDetail := _assistant.ToDetail()

	return assistantDetail, nil
}

func (s *assistantService) GetPaginatedList_ByCategoryAlias(
	ctx inner.Context,
	categoryAlias string,
	isPublicOnly bool,
	pageRequest pagination.PaginationRequest,
) (
	[]Assistant,
	pagination.PaginationMeta,
	error,
) {
	category, err := s.categoryService.GetOne_ByAlias(
		ctx,
		categoryAlias,
	)
	if err != nil {
		return nil, pagination.PaginationMeta{}, err
	}

	return s.assistantRepository.FindPaginatedList_ByCategoryID(
		ctx,
		category.ID,
		isPublicOnly,
		pageRequest,
	)
}

func (s *assistantService) GetPaginatedList_ByAuthor(
	ctx inner.Context,
	authorID string,
	pageRequest pagination.PaginationRequest,
) (
	[]Assistant,
	pagination.PaginationMeta,
	error,
) {
	return s.assistantRepository.FindPaginatedList_ByAuthorID(
		ctx,
		authorID,
		pageRequest,
	)
}

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

	_assister, err := s.assisterService.RegisterOne(
		ctx,
		assister.AssisterRegisterRequest{
			Origin:        assister.Origin_OpenAI,
			Model:         podoopenai.Model_GPT4oMini,
			Tests:         request.Tests,
			Fields:        request.Fields,
			QueryMessages: request.QueryMessages,
			Cost:          3,
		},
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return Assistant{}, err
	}

	_assistant, err := s.assistantRepository.InsertOne(
		ctx,
		request.ToModelForInsert(
			authorID,
			_assister.ID,
		),
	)
	if err != nil {
		s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		return Assistant{}, err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return Assistant{}, err
	}

	_assistant.Assister = _assister

	return _assistant, nil
}

func (s *assistantService) UpdateOne(
	ctx inner.Context,
	authorID string,
	request UpdateOneRequest,
) error {
	_assistant, err := s.assistantRepository.FindOne_ByID(
		ctx,
		request.ID,
		AssistantJoinOption{},
	)
	if err != nil {
		return err
	}
	if _assistant.AuthorID != authorID {
		return exception.ErrUnauthorized
	}

	_assister, err := s.assisterService.GetOne_ByID(
		ctx,
		_assistant.AssisterID,
	)
	if err != nil {
		return err
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

	newAssistant := request.ToModelForUpdate(_assistant)
	err = s.assistantRepository.UpdateOne(
		ctx,
		newAssistant,
	)
	if err != nil {
		s.cm.ClearTX(ctx, inner.TX_PodoSql)
		return err
	}

	newAssister := _assister.Copy()
	newAssister.Fields = request.Fields
	newAssister.Tests = request.Tests
	newAssister.QueryMessages = request.QueryMessages

	err = s.assisterService.UpdateOne(
		ctx,
		newAssister,
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
	_assistant, err := s.assistantRepository.FindOne_ByID(
		ctx,
		id,
		AssistantJoinOption{},
	)
	if err != nil {
		return err
	}
	if _assistant.Status != AssistantStatus_Registered || _assistant.AuthorID != authorID {
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

	if err := s.assisterService.RemoveOne_ByID(
		ctx,
		_assistant.AssisterID,
	); err != nil {
		s.cm.ClearTX(ctx, inner.TX_PodoSql)
		return err
	}

	if err := s.assistantRepository.DeleteOne_ByID(
		ctx,
		id,
	); err != nil {
		s.cm.ClearTX(ctx, inner.TX_PodoSql)
		return err
	}

	if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
		return err
	}

	return nil
}

func (s *assistantService) ApproveOne(
	ctx inner.Context,
	id string,
	cost int,
	metaTags []string,
) error {
	_assistant, err := s.assistantRepository.FindOne_ByID(
		ctx,
		id,
		AssistantJoinOption{},
	)
	if err != nil {
		return err
	}
	if _assistant.Status != AssistantStatus_UnderReview {
		return exception.ErrBadRequest
	}

	if cost > 0 {
		_assister, err := s.assisterService.GetOne_ByID(ctx, _assistant.AssisterID)
		if err != nil {
			return err
		}
		_assister.Cost = dt.UInt(cost)
		s.assisterService.UpdateOne(
			ctx,
			_assister,
		)
	}

	now := mydate.Now()
	_assistant.Status = AssistantStatus_Approved
	_assistant.MetaTags = metaTags
	_assistant.IsPublic = true
	_assistant.PublishedAt = &now
	err = s.assistantRepository.UpdateOne(
		ctx,
		_assistant,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewAssistantService(
	openaiClient *podoopenai.Client,
	assistantRepository AssistantRepository,
	categoryService category.CategoryService,
	assisterService assister.AssisterService,
	walletService wallet.WalletService,
	reviewService review.ReviewService,
	cm inner.ContextManager,
) AssistantService {
	return &assistantService{
		openaiClient:        openaiClient,
		assistantRepository: assistantRepository,
		categoryService:     categoryService,
		assisterService:     assisterService,
		walletService:       walletService,
		reviewService:       reviewService,
		cm:                  cm,
	}
}
