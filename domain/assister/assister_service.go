package assister

import (
	"fmt"
	"strings"

	"github.com/purplior/sbec/domain/ledger"
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/wallet"
	"github.com/purplior/sbec/infra/port/openai"
	"github.com/purplior/sbec/lib/dt"
)

type (
	AssisterService interface {
		GetOne_ByID(
			ctx inner.Context,
			id string,
		) (
			Assister,
			error,
		)

		RegisterOne(
			ctx inner.Context,
			request AssisterRegisterRequest,
		) (
			Assister,
			error,
		)

		UpdateOne(
			ctx inner.Context,
			assister Assister,
		) error

		RemoveOne_ByID(
			ctx inner.Context,
			id string,
		) error

		Request(
			ctx inner.Context,
			userId string,
			id string,
			inputs []AssisterInput,
		) (
			string,
			error,
		)

		RequestAsStream(
			ctx inner.Context,
			userId string,
			id string,
			inputs []AssisterInput,
			onInit func() error,
			onReceiveMessage func(msg string) error,
		) error
	}
)

type (
	assisterService struct {
		openaiClient       *openai.Client
		walletService      wallet.WalletService
		assisterRepository AssisterRepository
		cm                 inner.ContextManager
	}
)

func (s *assisterService) GetOne_ByID(
	ctx inner.Context,
	id string,
) (
	Assister,
	error,
) {
	return s.assisterRepository.FindOne_ByID(
		ctx,
		id,
	)
}

func (s *assisterService) RegisterOne(
	ctx inner.Context,
	request AssisterRegisterRequest,
) (
	Assister,
	error,
) {
	return s.assisterRepository.InsertOne(
		ctx,
		request.ToModelForInsert(),
	)
}

func (s *assisterService) UpdateOne(
	ctx inner.Context,
	assister Assister,
) error {
	return s.assisterRepository.UpdateOne(ctx, assister)
}

func (s *assisterService) RemoveOne_ByID(
	ctx inner.Context,
	id string,
) error {
	return s.assisterRepository.DeleteOne_ByID(
		ctx,
		id,
	)
}

func (s *assisterService) Request(
	ctx inner.Context,
	userId string,
	id string,
	inputs []AssisterInput,
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

	messageData, err := s.createMessageData(assister, inputs)
	if err != nil {
		return "", err
	}
	messages, err := s.createMessagesOfGPT(assister, messageData)
	if err != nil {
		return "", err
	}

	if !isFree {
		if err := s.cm.BeginTX(ctx, inner.TX_sqldb); err != nil {
			return "", err
		}

		if err := s.walletService.Expend(
			ctx,
			userId,
			int(assister.Cost),
			ledger.LedgerAction_ConsumeByAssister,
			assister.ID,
		); err != nil {
			s.cm.RollbackTX(ctx, inner.TX_sqldb)
			return "", err
		}
	}

	return s.openaiClient.RequestChatCompletions(
		ctx.Value(),
		openai.ChatCompletionRequest{
			Model:    string(assister.Model),
			Messages: messages,
		},
		func() error {
			if !isFree {
				if err := s.cm.CommitTX(ctx, inner.TX_sqldb); err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func (s *assisterService) RequestAsStream(
	ctx inner.Context,
	userId string,
	id string,
	inputs []AssisterInput,
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

	messageData, err := s.createMessageData(assister, inputs)
	if err != nil {
		return err
	}
	messages, err := s.createMessagesOfGPT(assister, messageData)
	if err != nil {
		return err
	}

	if !isFree {
		if err := s.cm.BeginTX(ctx, inner.TX_sqldb); err != nil {
			return err
		}

		if err := s.walletService.Expend(
			ctx,
			userId,
			int(assister.Cost),
			ledger.LedgerAction_ConsumeByAssister,
			assister.ID,
		); err != nil {
			s.cm.RollbackTX(ctx, inner.TX_sqldb)
			return err
		}
	}

	err = s.openaiClient.RequestChatCompletionsStream(
		ctx.Value(),
		openai.ChatCompletionRequest{
			Model:    string(assister.Model),
			Messages: messages,
		},
		func() error {
			if err := onInit(); err != nil {
				return err
			}
			if !isFree {
				if err := s.cm.CommitTX(ctx, inner.TX_sqldb); err != nil {
					return err
				}
			}

			return nil
		},
		onReceiveMessage,
	)
	if err != nil {
		if !isFree {
			s.cm.RollbackTX(ctx, inner.TX_sqldb)
		}
	}

	return err
}

func (s *assisterService) createMessageData(
	_assister Assister,
	inputs []AssisterInput,
) (map[string]string, error) {
	fieldTypeMap := map[string]AssisterFieldType{}
	for _, field := range _assister.Fields {
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
		case AssisterFieldType_Keyword:
			for i, value := range input.Values {
				if i > 0 {
					values += ","
				}
				values += value.(string)
			}
		case AssisterFieldType_Paragraph:
			for i, value := range input.Values {
				if i > 0 {
					values += "\n"
				}
				values += "- " + value.(string)
			}
		case AssisterFieldType_ParagraphGroup:
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
	_assister Assister,
	data map[string]string,
) ([]map[string]string, error) {
	messages := make([]map[string]string, len(_assister.QueryMessages))

	for i, queryMessage := range _assister.QueryMessages {
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
	openaiClient *openai.Client,
	walletService wallet.WalletService,
	assisterRepository AssisterRepository,
	cm inner.ContextManager,
) AssisterService {
	return &assisterService{
		openaiClient:       openaiClient,
		walletService:      walletService,
		assisterRepository: assisterRepository,
		cm:                 cm,
	}
}
