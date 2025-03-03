package assister

import (
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/domain/wallet"
	"github.com/purplior/edi-adam/domain/walletlog"
	"github.com/purplior/edi-adam/infra/port/openai"
)

type (
	AssisterService interface {
		Get(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Assister,
			error,
		)

		Register(
			session inner.Session,
			dto RegisterDTO,
		) (
			model.Assister,
			error,
		)

		Update(
			session inner.Session,
			dto UpdateDTO,
		) error

		Remove(
			session inner.Session,
			queryOption QueryOption,
		) error

		Request(
			session inner.Session,
			dto RequestDTO,
		) (
			string,
			error,
		)

		RequestAsStream(
			session inner.Session,
			dto RequestDTO,
			onInit func() error,
			onReceiveMessage func(msg string) error,
		) error
	}
)

type (
	assisterService struct {
		openAI             *assisterOpenAI
		walletService      wallet.WalletService
		assisterRepository AssisterRepository
	}
)

func (s *assisterService) Get(
	session inner.Session,
	queryOption QueryOption,
) (
	model.Assister,
	error,
) {
	return s.assisterRepository.Read(
		session,
		queryOption,
	)
}

func (s *assisterService) Register(
	session inner.Session,
	dto RegisterDTO,
) (
	model.Assister,
	error,
) {
	return s.assisterRepository.Create(
		session,
		model.Assister{
			ID:            dto.ID,
			Origin:        dto.Origin,
			Model:         dto.Model,
			Fields:        dto.Fields,
			QueryMessages: dto.QueryMessages,
			Cost:          3,
			Tests:         dto.Tests,
			AssistantID:   dto.AssistantID,
		},
	)
}

func (s *assisterService) Update(
	session inner.Session,
	dto UpdateDTO,
) error {
	return s.assisterRepository.Updates(
		session,
		QueryOption{
			ID: dto.ID,
		},
		model.Assister{
			ID:            dto.ID,
			Origin:        dto.Origin,
			Model:         dto.Model,
			Fields:        dto.Fields,
			QueryMessages: dto.QueryMessages,
			Tests:         dto.Tests,
		},
	)
}

func (s *assisterService) Remove(
	session inner.Session,
	queryOption QueryOption,
) error {
	return s.assisterRepository.Delete(
		session,
		queryOption,
	)
}

func (s *assisterService) Request(
	session inner.Session,
	dto RequestDTO,
) (
	result string,
	err error,
) {
	var m model.Assister
	if m, err = s.assisterRepository.Read(
		session,
		QueryOption{
			ID: dto.ID,
		},
	); err != nil {
		return "", err
	}
	// 로그인 하지 않은 사용자는 Cost
	if dto.UserID == 0 && m.Cost > 0 {
		return "", exception.ErrBadRequest
	}
	if m.Cost > 0 {
		if err = s.walletService.Expend(
			session,
			wallet.ExpendDTO{
				OwnerID: dto.UserID,
				Delta:   m.Cost,
				LogAddDTO: walletlog.AddDTO{
					Type:    model.WalletLogType_ExpendOnUsingAssister,
					Delta:   int(m.Cost),
					Comment: "샘비서 ID: " + m.ID,
				},
			},
		); err != nil {
			return "", err
		}
	}

	switch m.Origin {
	case model.AssisterOrigin_OpenAI:
		result, err = s.openAI.RequestChatCompletions(
			session,
			m,
			dto.Inputs,
		)
	}

	return result, err
}

func (s *assisterService) RequestAsStream(
	session inner.Session,
	dto RequestDTO,
	onInit func() error,
	onReceiveMessage func(msg string) error,
) (err error) {
	var m model.Assister
	if m, err = s.assisterRepository.Read(
		session,
		QueryOption{
			ID: dto.ID,
		},
	); err != nil {
		return err
	}

	if dto.UserID == 0 && m.Cost > 0 {
		return exception.ErrBadRequest
	}
	if m.Cost > 0 {
		if err = s.walletService.Expend(
			session,
			wallet.ExpendDTO{
				OwnerID: dto.UserID,
				Delta:   m.Cost,
				LogAddDTO: walletlog.AddDTO{
					Type:    model.WalletLogType_ExpendOnUsingAssister,
					Delta:   int(m.Cost),
					Comment: "샘비서 ID: " + m.ID,
				},
			},
		); err != nil {
			return err
		}
	}

	switch m.Origin {
	case model.AssisterOrigin_OpenAI:
		err = s.openAI.RequestChatCompletionsAsStream(
			session,
			m,
			dto.Inputs,
			onInit,
			onReceiveMessage,
		)
	}

	return err
}

func NewAssisterService(
	openaiClient *openai.Client,
	walletService wallet.WalletService,
	assisterRepository AssisterRepository,
) AssisterService {
	return &assisterService{
		openAI: &assisterOpenAI{
			openaiClient: openaiClient,
		},
		walletService:      walletService,
		assisterRepository: assisterRepository,
	}
}
