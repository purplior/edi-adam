package assister

import (
	"fmt"
	"strings"

	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/port/openai"
	"github.com/purplior/edi-adam/lib/dt"
)

type (
	assisterOpenAI struct {
		openaiClient *openai.Client
	}
)

func (s *assisterOpenAI) RequestChatCompletions(
	session inner.Session,
	m model.Assister,
	inputs []model.AssisterInput,
) (string, error) {
	messageData, err := s.createMessageData(m, inputs)
	if err != nil {
		return "", err
	}
	messages, err := s.createMessagesOfGPT(m, messageData)
	if err != nil {
		return "", err
	}

	return s.openaiClient.RequestChatCompletions(
		session.Context(),
		openai.ChatCompletionRequest{
			Model:    string(m.Model),
			Messages: messages,
		},
	)
}

func (s *assisterOpenAI) RequestChatCompletionsAsStream(
	session inner.Session,
	m model.Assister,
	inputs []model.AssisterInput,
	onInit func() error,
	onReceiveMessage func(message string) error,
) error {
	messageData, err := s.createMessageData(m, inputs)
	if err != nil {
		return err
	}
	messages, err := s.createMessagesOfGPT(m, messageData)
	if err != nil {
		return err
	}

	err = s.openaiClient.RequestChatCompletionsStream(
		session.Context(),
		openai.ChatCompletionRequest{
			Model:       string(m.Model),
			Messages:    messages,
			Temperature: 0.7,
			TopP:        0.8,
		},
		onInit,
		onReceiveMessage,
	)

	return err
}

func (s *assisterOpenAI) createMessageData(
	_assister model.Assister,
	inputs []model.AssisterInput,
) (map[string]string, error) {
	fieldTypeMap := map[string]model.AssisterFieldType{}
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
		case model.AssisterFieldType_Keyword:
			for i, value := range input.Values {
				if i > 0 {
					values += ","
				}
				values += value.(string)
			}
		case model.AssisterFieldType_Paragraph:
			for i, value := range input.Values {
				if i > 0 {
					values += "\n"
				}
				values += "- " + value.(string)
			}
		case model.AssisterFieldType_ParagraphGroup:
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

func (s *assisterOpenAI) createMessagesOfGPT(
	_assister model.Assister,
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
