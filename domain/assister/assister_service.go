package assister

import (
	"log"

	"github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/infra/port/podoopenai"
)

type (
	AssisterService interface {
		RequestStream(
			ctx context.APIContext,
			id string,
			inputs map[string]interface{},
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
	inputs map[string]interface{},
	onInit func() error,
	onReceiveMessage func(msg string) error,
) error {
	assisterForm, err := s.assisterFormService.GetOneByAssisterID(ctx, id)
	if err != nil {
		return err
	}

	// 요청문장 가공
	log.Println(assisterForm.CreatedAt)
	messages := []map[string]string{
		{"role": "user", "content": "선생님들이 겪는 어려움에 대해 1000자로 정리해줘"},
	}

	return s.openaiClient.RequestChatCompletionsStream(
		ctx,
		podoopenai.ChatCompletionRequest{
			Model:    string(AssisterModel_OpenAI_ChatGPT4o),
			Messages: messages,
		},
		onInit,
		onReceiveMessage,
	)
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
