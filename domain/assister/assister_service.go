package assister

import (
	"github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/domain/shared/exception"
	"github.com/podossaem/podoroot/infra/port/podoopenai"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	AssisterService interface {
		RequestStream(
			ctx context.APIContext,
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
	}
)

func (s *assisterService) RequestStream(
	ctx context.APIContext,
	id string,
	inputs []assisterform.AssisterInput,
	onInit func() error,
	onReceiveMessage func(msg string) error,
) error {
	form, err := s.assisterFormService.GetOneByAssisterID(ctx, id)
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

	print("info:\n")
	print(info)
	for i, message := range form.QueryMessages {
		if i == messageLen-1 {
			message.Content += "\n\n" + info
		}
		messages[i] = message.CreatePayload()
	}

	return s.openaiClient.RequestChatCompletionsStream(
		ctx,
		podoopenai.ChatCompletionRequest{
			Model:    string(form.Model),
			Messages: messages,
		},
		onInit,
		onReceiveMessage,
	)
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
) AssisterService {
	return &assisterService{
		openaiClient:        openaiClient,
		assisterFormService: assisterFormService,
	}
}
