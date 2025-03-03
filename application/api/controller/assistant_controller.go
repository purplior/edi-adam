package controller

import (
	"net/http"

	"github.com/purplior/edi-adam/application/common"
	domain "github.com/purplior/edi-adam/domain/assistant"
	"github.com/purplior/edi-adam/domain/assister"
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/lib/strgen"
)

type (
	AssistantController interface {
		Controller
		// 샘어시 상세정보 가져오기
		Get_Detail() common.Route
		// 샘어시 목록 가져오기 (카테고리)
		GetPaginatedList_Card() common.Route
		// 내가 만든 샘어시 목록 가져오기
		GetMyPaginatedList_Card() common.Route
		// 나의 새로운 샘어시 등록하기
		RegisterMine() common.Route
		// 나의 샘어시 수정하기
		UpdateMine() common.Route
		// 나의 샘어시 제거하기
		RemoveMine() common.Route
	}
)

type (
	assistantController struct {
		assistantService domain.AssistantService
		assisterService  assister.AssisterService
	}
)

func (c *assistantController) GroupPath() string {
	return "/assistants"
}

func (c *assistantController) Routes() []common.Route {
	return []common.Route{
		c.Get_Detail(),
		c.GetPaginatedList_Card(),
		c.GetMyPaginatedList_Card(),
		c.RegisterMine(),
		c.UpdateMine(),
		c.RemoveMine(),
	}
}

func (c *assistantController) Get_Detail() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/o/:id/detail",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			ID uint `param:"id,required"`
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		_assistant, err := c.assistantService.Get(
			ctx.Session(),
			domain.QueryOption{
				ID:           dto.ID,
				IsPublic:     true,
				WithCategory: true,
				WithAuthor:   true,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}
		if !_assistant.IsPublic {
			return ctx.SendError(exception.ErrUnauthorized)
		}

		_assister, err := c.assisterService.Get(
			ctx.Session(),
			assister.QueryOption{
				ID: _assistant.CurrentAssisterID,
			},
		)
		if err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(struct {
			Assistant model.Assistant `json:"assistant"`
			Assister  model.Assister  `json:"assister"`
		}{
			Assistant: _assistant,
			Assister:  _assister,
		})
	}

	return route
}

func (c *assistantController) GetPaginatedList_Card() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Path:   "/o/lst/card",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			CategoryID string `query:"cid,required"`
			pagination.PaginationRequest
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		return c.sendCard_PaginatedList(
			ctx,
			pagination.PaginationQuery[domain.QueryOption]{
				QueryOption: domain.QueryOption{
					CategoryID:   dto.CategoryID,
					IsPublic:     true,
					WithCategory: true,
					WithAuthor:   true,
				},
				PageRequest: dto.PaginationRequest,
			},
		)
	}

	return route
}

func (c *assistantController) GetMyPaginatedList_Card() common.Route {
	route := common.Route{
		Method: http.MethodGet,
		Option: common.RouteOption{Member: true},
		Path:   "/m/lst/card",
	}
	route.Handler = func(ctx *common.Context) error {
		var dto struct {
			pagination.PaginationRequest
		}
		if err := ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		return c.sendCard_PaginatedList(
			ctx,
			pagination.PaginationQuery[domain.QueryOption]{
				QueryOption: domain.QueryOption{
					AuthorID:     ctx.Identity.ID,
					IsPublic:     true,
					WithCategory: true,
					WithAuthor:   true,
				},
				PageRequest: dto.PaginationRequest,
			},
		)
	}

	return route
}

func (c *assistantController) RegisterMine() common.Route {
	route := common.Route{
		Method: http.MethodPost,
		Option: common.RouteOption{Member: true},
		Path:   "/m/fst",
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			Icon        string   `body:"icon"`
			Title       string   `body:"title"`
			Description string   `body:"description"`
			Notice      string   `body:"notice"`
			CategoryID  string   `body:"categoryId"`
			Tags        []string `body:"tags"`
			IsPublic    bool     `body:"isPublic"`

			Origin        string                       `body:"origin"`
			Model         string                       `body:"model"`
			Fields        []model.AssisterField        `body:"fields"`
			QueryMessages []model.AssisterQueryMessage `body:"queryMessages"`
			Tests         []model.AssisterInput        `body:"tests"`
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		assisterID, err := strgen.UniqueSortableID()
		if err != nil {
			return ctx.SendError(err)
		}

		session := ctx.Session()
		if err = session.BeginTransaction(); err != nil {
			return ctx.SendError(err)
		}

		assistantRegisterDTO := domain.RegisterDTO{
			UserID:            ctx.Identity.ID,
			CurrentAssisterID: assisterID,

			Icon:        dto.Icon,
			Title:       dto.Title,
			Description: dto.Description,
			Notice:      dto.Notice,
			CategoryID:  dto.CategoryID,
			Tags:        dto.Tags,
			IsPublic:    dto.IsPublic,
		}
		if !assistantRegisterDTO.IsValid() {
			return ctx.SendError(exception.ErrBadRequest)
		}
		assisterRegisterDTO := assister.RegisterDTO{
			ID:            assisterID,
			Origin:        model.AssisterOrigin(dto.Origin),
			Model:         model.AssisterModel(dto.Model),
			Fields:        dto.Fields,
			QueryMessages: dto.QueryMessages,
			Tests:         dto.Tests,
		}
		if !assisterRegisterDTO.IsValid() {
			return ctx.SendError(exception.ErrBadRequest)
		}

		_assistant, err := c.assistantService.Register(
			session,
			assistantRegisterDTO,
		)
		if err != nil {
			session.RollbackTransaction()
			return ctx.SendError(err)
		}

		assisterRegisterDTO.AssistantID = _assistant.ID
		_, err = c.assisterService.Register(
			session,
			assisterRegisterDTO,
		)
		if err != nil {
			session.RollbackTransaction()
			return ctx.SendError(err)
		}

		if err = session.CommitTransaction(); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendCreated(nil)
	}

	return route
}

func (c *assistantController) UpdateMine() common.Route {
	route := common.Route{
		Method: http.MethodPatch,
		Option: common.RouteOption{Member: true},
		Path:   "/m/:id",
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			ID uint `param:"id,required"`

			Icon        string   `body:"icon"`
			Title       string   `body:"title"`
			Description string   `body:"description"`
			Notice      string   `body:"notice"`
			CategoryID  string   `body:"categoryId"`
			Tags        []string `body:"tags"`
			IsPublic    bool     `body:"isPublic"`

			AssisterID    string                       `body:"assisterId"`
			Origin        string                       `body:"origin"`
			Model         string                       `body:"model"`
			Fields        []model.AssisterField        `body:"fields"`
			QueryMessages []model.AssisterQueryMessage `body:"queryMessages"`
			Tests         []model.AssisterInput        `body:"tests"`
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		if err = c.assistantService.Update(
			ctx.Session(),
			domain.UpdateDTO{
				ID:          dto.ID,
				UserID:      ctx.Identity.ID,
				Icon:        dto.Icon,
				Title:       dto.Title,
				Description: dto.Description,
				Notice:      dto.Notice,
				CategoryID:  dto.CategoryID,
				Tags:        dto.Tags,
				IsPublic:    dto.IsPublic,
			},
		); err != nil {
			return ctx.SendError(err)
		}

		if err = c.assisterService.Update(
			ctx.Session(),
			assister.UpdateDTO{
				ID:            dto.AssisterID,
				Origin:        model.AssisterOrigin(dto.Origin),
				Model:         model.AssisterModel(dto.Model),
				Fields:        dto.Fields,
				QueryMessages: dto.QueryMessages,
				Tests:         dto.Tests,
			},
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(nil)
	}

	return route
}

func (c *assistantController) RemoveMine() common.Route {
	route := common.Route{
		Method: http.MethodPatch,
		Option: common.RouteOption{Member: true},
		Path:   "/m/:id",
	}
	route.Handler = func(ctx *common.Context) (err error) {
		var dto struct {
			ID uint `param:"id,required"`
		}
		if err = ctx.BindRequestData(&dto); err != nil {
			return ctx.SendError(err)
		}

		if err = c.assistantService.Remove(
			ctx.Session(),
			domain.RemoveDTO{
				UserID: ctx.Identity.ID,
				ID:     dto.ID,
			},
		); err != nil {
			return ctx.SendError(err)
		}

		return ctx.SendOK(nil)
	}

	return route
}

func (c *assistantController) sendCard_PaginatedList(
	ctx *common.Context,
	query pagination.PaginationQuery[domain.QueryOption],
) error {
	assistants, pageMeta, err := c.assistantService.GetPaginatedList(
		ctx.Session(),
		query,
	)
	if err != nil && err != exception.ErrNoRecord {
		return ctx.SendError(err)
	}

	assistantCards := make([]model.AssistantCard, len(assistants))
	for i, assistant := range assistants {
		assistantCards[i] = assistant.ToCard()
	}

	return ctx.SendOK(struct {
		AssistantCards []model.AssistantCard     `json:"assistantCards"`
		PageMeta       pagination.PaginationMeta `json:"pageMeta"`
	}{
		AssistantCards: assistantCards,
		PageMeta:       pageMeta,
	})
}

func NewAssistantController(
	assistantService domain.AssistantService,
	assisterService assister.AssisterService,
) AssistantController {
	return &assistantController{
		assistantService: assistantService,
		assisterService:  assisterService,
	}
}
