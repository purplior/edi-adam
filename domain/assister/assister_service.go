package assister

import (
	"github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/ledger"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/wallet"
	"github.com/podossaem/podoroot/infra/port/podoopenai"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	AssisterService interface {
		RequestStream(
			ctx inner.Context,
			userId string,
			id string,
			inputs []assisterform.AssisterInput,
			onInit func() error,
			onReceiveMessage func(msg string) error,
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

	info, err := s.createQueryInformation(
		form,
		inputs,
	)
	if err != nil {
		return err
	}

	messageLen := len(form.QueryMessages)
	messages := make([]map[string]string, messageLen)

	for i, message := range form.QueryMessages {
		if i == messageLen-1 {
			message.Content += "\n\n" + info
		}
		messages[i] = message.CreatePayload()
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
			assister.ViewID,
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

func (s *assisterService) createQueryInformation(
	form assisterform.AssisterForm,
	inputs []assisterform.AssisterInput,
) (string, error) {
	content := "## "
	if len(form.QueryInfoHeading) == 0 {
		content += "정보"
	} else {
		content += form.QueryInfoHeading
	}

	for _, input := range inputs {
		content += "\n\n### " + input.Name

		for i, v := range input.Values {
			field, ok := form.FindField(input.Name)
			if !ok {
				continue
			}

			switch field.Type {
			case assisterform.AssisterFieldType_Keyword:
				{
					value := v.(string)
					if len(value) == 0 {
						if field.Required {
							return "", exception.ErrBadRequest
						} else {
							continue
						}
					}

					if i > 0 {
						content += ","
					} else {
						content += "\n"
					}
					content += value
				}
			case assisterform.AssisterFieldType_Paragraph:
				{
					value := v.(string)
					if len(value) == 0 {
						if field.Required {
							return "", exception.ErrBadRequest
						} else {
							continue
						}
					}

					content += "\n" + dt.Str(i+1) + ". " + value
				}
			case assisterform.AssisterFieldType_ParagraphGroup:
				{
					vIObjects := v.([]interface{})
					vIObjectsLen := len(vIObjects)
					if vIObjectsLen == 0 {
						if field.Required {
							return "", exception.ErrBadRequest
						} else {
							continue
						}
					}

					content += "\n" + dt.Str(i+1) + ". " + field.ItemName
					for _, vIObj := range vIObjects {
						vObj := vIObj.(map[string]interface{})
						childName := dt.Str(vObj["name"])
						childValue := dt.Str(vObj["value"])

						if len(childValue) == 0 {
							if field.Required {
								return "", exception.ErrBadRequest
							} else {
								continue
							}
						}

						content += "\n\t- " + childName + ": " + childValue
					}
				}
			}
		}
	}

	return content, nil
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
