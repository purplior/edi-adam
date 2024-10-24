package assister

import (
	"github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/context"
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
	assisterForm, err := s.assisterFormService.GetOneByAssisterID(ctx, id)
	if err != nil {
		return err
	}

	info := s.createQueryInformation(assisterForm.QueryInfoHeading, inputs)
	messageLen := len(assisterForm.QueryMessages)
	messages := make([]map[string]string, messageLen)

	print("info:\n")
	print(info)
	for i, message := range assisterForm.QueryMessages {
		if i == messageLen-1 {
			message.Content += "\n\n" + info
		}
		messages[i] = message.CreatePayload()
	}

	return s.openaiClient.RequestChatCompletionsStream(
		ctx,
		podoopenai.ChatCompletionRequest{
			Model:    string(assisterForm.Model),
			Messages: messages,
		},
		onInit,
		onReceiveMessage,
	)
}

func (s *assisterService) createQueryInformation(
	queryInfoHeading string,
	inputs []assisterform.AssisterInput,
) string {
	content := "## "
	if len(queryInfoHeading) == 0 {
		content += "정보"
	} else {
		content += queryInfoHeading
	}

	for _, input := range inputs {
		content += "\n\n### " + input.Name

		for i, v := range input.Values {
			switch input.Type {
			case assisterform.AssisterFieldType_Keyword:
				{
					value := v.(string)
					if len(value) == 0 {
						continue
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
						continue
					}

					content += "\n" + dt.Str(i+1) + ". " + value
				}
			case assisterform.AssisterFieldType_ParagraphGroup:
				{
					vIObjects := v.([]interface{})
					vIObjectsLen := len(vIObjects)
					if vIObjectsLen == 0 {
						continue
					}

					content += "\n" + dt.Str(i+1) + ". " + input.ItemName
					for _, vIObj := range vIObjects {
						vObj := vIObj.(map[string]interface{})
						childName := dt.Str(vObj["name"])
						childValue := dt.Str(vObj["value"])

						if len(childValue) == 0 {
							continue
						}

						content += "\n\t- " + childName + ": " + childValue
					}
				}
			}
		}

	}

	return content
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
