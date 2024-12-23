package assister

import (
	"fmt"
	"strings"

	"github.com/purplior/podoroot/domain/assisterform"
	"github.com/purplior/podoroot/domain/ledger"
	"github.com/purplior/podoroot/domain/shared/exception"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/domain/wallet"
	"github.com/purplior/podoroot/infra/port/podoopenai"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	AssisterService interface {
		RegisterOne(
			ctx inner.Context,
			request AssisterRegisterRequest,
		) (
			Assister,
			error,
		)

		Request(
			ctx inner.Context,
			userId string,
			id string,
			inputs []assisterform.AssisterInput,
		) (
			string,
			error,
		)

		RequestStream(
			ctx inner.Context,
			userId string,
			id string,
			inputs []assisterform.AssisterInput,
			onInit func() error,
			onReceiveMessage func(msg string) error,
		) error

		GetOne_ByID(
			ctx inner.Context,
			id string,
		) (
			Assister,
			error,
		)

		GetPaginatedList_ByAssistant(
			ctx inner.Context,
			assistantID string,
			page int,
			pageSize int,
		) (
			[]Assister,
			pagination.PaginationMeta,
			error,
		)

		PutOne(
			ctx inner.Context,
			assister Assister,
		) error

		CreateOne(
			ctx inner.Context,
			assister Assister,
		) error

		RemoveAll_ByIDs(
			ctx inner.Context,
			ids []string,
		) error
	}
)

type (
	assisterService struct {
		openaiClient        *podoopenai.Client
		assisterFormService assisterform.AssisterFormService
		walletService       wallet.WalletService
		assisterRepository  AssisterRepository
		cm                  inner.ContextManager
	}
)

func (s *assisterService) RegisterOne(
	ctx inner.Context,
	request AssisterRegisterRequest,
) (
	Assister,
	error,
) {
	return s.assisterRepository.InsertOne(
		ctx,
		Assister{
			AssistantID:        request.AssistantID,
			Version:            request.Version,
			VersionDescription: request.VersionDescription,
			Cost:               request.Cost,
		},
	)
}

func (s *assisterService) Request(
	ctx inner.Context,
	userId string,
	id string,
	inputs []assisterform.AssisterInput,
) (
	string,
	error,
) {
	assister, err := s.assisterRepository.FindOne_ByID(ctx, id)
	if err != nil {
		return "", err
	}

	isFree := assister.Cost == 0
	if !isFree && len(userId) == 0 {
		return "", exception.ErrBadRequest
	}

	form, err := s.assisterFormService.GetOne_ByAssisterID(ctx, id)
	if err != nil {
		return "", err
	}

	messageData, err := s.createMessageData(form, inputs)
	if err != nil {
		return "", err
	}
	messages, err := s.createMessagesOfGPT(form, messageData)
	if err != nil {
		return "", err
	}

	if !isFree {
		if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
			return "", err
		}

		if err := s.walletService.Expend(
			ctx,
			userId,
			int(assister.Cost),
			ledger.LedgerAction_ConsumeByAssister,
			assister.ID,
		); err != nil {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
			return "", err
		}
	}

	return s.openaiClient.RequestChatCompletions(
		ctx.Value(),
		podoopenai.ChatCompletionRequest{
			Model:    string(form.Model),
			Messages: messages,
		},
		func() error {
			if !isFree {
				if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func (s *assisterService) RequestStream(
	ctx inner.Context,
	userId string,
	id string,
	inputs []assisterform.AssisterInput,
	onInit func() error,
	onReceiveMessage func(msg string) error,
) error {
	assister, err := s.assisterRepository.FindOne_ByID(ctx, id)
	if err != nil {
		return err
	}

	isFree := assister.Cost == 0
	if !isFree && len(userId) == 0 {
		return exception.ErrBadRequest
	}

	form, err := s.assisterFormService.GetOne_ByAssisterID(ctx, id)
	if err != nil {
		return err
	}

	messageData, err := s.createMessageData(form, inputs)
	if err != nil {
		return err
	}
	messages, err := s.createMessagesOfGPT(form, messageData)
	if err != nil {
		return err
	}

	if !isFree {
		if err := s.cm.BeginTX(ctx, inner.TX_PodoSql); err != nil {
			return err
		}

		if err := s.walletService.Expend(
			ctx,
			userId,
			int(assister.Cost),
			ledger.LedgerAction_ConsumeByAssister,
			assister.ID,
		); err != nil {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
			return err
		}
	}

	err = s.openaiClient.RequestChatCompletionsStream(
		ctx.Value(),
		podoopenai.ChatCompletionRequest{
			Model:    string(form.Model),
			Messages: messages,
		},
		func() error {
			if err := onInit(); err != nil {
				return err
			}
			if !isFree {
				if err := s.cm.CommitTX(ctx, inner.TX_PodoSql); err != nil {
					return err
				}
			}

			return nil
		},
		onReceiveMessage,
	)
	if err != nil {
		if !isFree {
			s.cm.RollbackTX(ctx, inner.TX_PodoSql)
		}
	}

	return err
}

func (s *assisterService) GetOne_ByID(
	ctx inner.Context,
	id string,
) (
	Assister,
	error,
) {
	return s.assisterRepository.FindOne_ByID(ctx, id)
}

func (s *assisterService) GetPaginatedList_ByAssistant(
	ctx inner.Context,
	assistantID string,
	page int,
	pageSize int,
) (
	[]Assister,
	pagination.PaginationMeta,
	error,
) {
	return s.assisterRepository.FindPaginatedList_ByAssistantID(
		ctx,
		assistantID,
		page,
		pageSize,
	)
}

func (s *assisterService) PutOne(
	ctx inner.Context,
	assister Assister,
) error {
	return s.assisterRepository.UpdateOne(ctx, assister)
}

func (s *assisterService) CreateOne(
	ctx inner.Context,
	assister Assister,
) error {
	_, err := s.assisterRepository.InsertOne(
		ctx,
		assister,
	)

	return err
}

func (s *assisterService) RemoveAll_ByIDs(
	ctx inner.Context,
	ids []string,
) error {
	return s.assisterRepository.DeleteAll_ByIDs(ctx, ids)
}

func (s *assisterService) createMessageData(
	form assisterform.AssisterForm,
	inputs []assisterform.AssisterInput,
) (map[string]string, error) {
	fieldTypeMap := map[string]assisterform.AssisterFieldType{}
	for _, field := range form.Fields {
		fieldTypeMap[field.Name] = field.Type
	}

	data := map[string]string{}
	for _, input := range inputs {
		fieldType, ok := fieldTypeMap[input.Name]
		if !ok {
			return nil, ErrInvalidAssisterInput
		}

		values := ""

		switch fieldType {
		case assisterform.AssisterFieldType_Keyword:
			for i, value := range input.Values {
				if i > 0 {
					values += ","
				}
				values += value.(string)
			}
		case assisterform.AssisterFieldType_Paragraph:
			for i, value := range input.Values {
				if i > 0 {
					values += "\n"
				}
				values += "- " + value.(string)
			}
		case assisterform.AssisterFieldType_ParagraphGroup:
			for i, value := range input.Values {
				if i > 0 {
					values += "\n"
				}
				values += dt.Str(i+1) + ". "

				childrenInterface := value.([]interface{})
				for _, childInterface := range childrenInterface {
					child := childInterface.(map[string]interface{})
					childName := child["name"].(string)
					childValue := child["value"].(string)

					values += "\n\t- " + childName + ": " + childValue
				}
			}
		}

		data[input.Name] = values
	}

	return data, nil
}

func (s *assisterService) createMessagesOfGPT(
	form assisterform.AssisterForm,
	data map[string]string,
) ([]map[string]string, error) {
	messages := make([]map[string]string, len(form.QueryMessages))

	for i, queryMessage := range form.QueryMessages {
		template := queryMessage.Content
		for key, val := range data {
			placeholder := fmt.Sprintf("{{ %s }}", key)
			template = strings.ReplaceAll(template, placeholder, val)
		}

		messages[i] = map[string]string{
			"role":    string(queryMessage.Role),
			"content": template,
		}
	}

	return messages, nil
}

func NewAssisterService(
	openaiClient *podoopenai.Client,
	assisterFormService assisterform.AssisterFormService,
	walletService wallet.WalletService,
	assisterRepository AssisterRepository,
	cm inner.ContextManager,
) AssisterService {
	return &assisterService{
		openaiClient:        openaiClient,
		assisterFormService: assisterFormService,
		walletService:       walletService,
		assisterRepository:  assisterRepository,
		cm:                  cm,
	}
}
