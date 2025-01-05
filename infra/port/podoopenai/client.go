package podoopenai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/domain/shared/constant"
	"github.com/purplior/podoroot/domain/shared/logger"
)

const (
	Model_GPT4o     = "chatgpt-4o-latest"
	Model_GPT4oMini = "gpt-4o-mini"

	Token_1M                              = 1000 * 1000
	GPT4o_Input_CostPerToken      float64 = (2.5 * constant.ExchangeWonOfDollar) / Token_1M
	GPT4o_Output_CostPerToken     float64 = (10 * constant.ExchangeWonOfDollar) / Token_1M
	GPT4oMini_Input_CostPerToken  float64 = (1.0 * constant.ExchangeWonOfDollar) / Token_1M
	GPT4oMini_Output_CostPerToken float64 = (0.6 * constant.ExchangeWonOfDollar) / Token_1M
)

type (
	ConstructorOption struct {
		ApiKey         string
		OrganizationID string
		ProjectID      string
	}

	Client struct {
		opt ConstructorOption
	}
)

type (
	ChatCompletionRequest struct {
		Model       string
		Messages    []map[string]string
		Temperature float64
		TopP        float64
	}

	chatCompletionResponseChunk struct {
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}
)

const (
	BufferUnit = 20
)

func (c *Client) RequestChatCompletions(
	ctx context.Context,
	request ChatCompletionRequest,
	onInit func() error,
) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       request.Model,
		"messages":    request.Messages,
		"temperature": request.Temperature,
		"top_p":       request.TopP,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+c.opt.ApiKey)
	req.Header.Set("OpenAI-Organization", c.opt.OrganizationID)
	req.Header.Set("OpenAI-Project", c.opt.ProjectID)

	httpClient := &http.Client{
		Timeout: time.Minute * 5,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	if err := onInit(); err != nil {
		return "", err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var completionResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens       int `json:"prompt_tokens"`
			CompletionTokens   int `json:"completion_tokens"`
			TotalTokens        int `json:"total_tokens"`
			PromptTokensDetail struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(responseBody, &completionResponse); err != nil {
		return "", err
	}

	if len(completionResponse.Choices) == 0 {
		return "", fmt.Errorf("no completion choices found")
	}

	content := completionResponse.Choices[0].Message.Content

	if config.DebugMode() {
		promptTokens := completionResponse.Usage.PromptTokens
		promptCachedTokens := completionResponse.Usage.PromptTokensDetail.CachedTokens
		completionTokens := completionResponse.Usage.CompletionTokens
		totalTokens := completionResponse.Usage.TotalTokens

		logger.Debug("prompt_token: %d", promptTokens)
		logger.Debug("cached_tokens: %d", promptCachedTokens)
		logger.Debug("completion_tokens: %d", completionTokens)
		logger.Debug("total_tokens: %d", totalTokens)

		inputCostPerToken := 1.0
		outputCostPerToken := 1.0
		switch request.Model {
		case Model_GPT4o:
			inputCostPerToken = GPT4o_Input_CostPerToken
			outputCostPerToken = GPT4o_Output_CostPerToken
		case Model_GPT4oMini:
			inputCostPerToken = GPT4oMini_Input_CostPerToken
			outputCostPerToken = GPT4oMini_Output_CostPerToken
		}
		prompt_won := inputCostPerToken * float64(promptTokens)
		completion_won := outputCostPerToken * float64(completionTokens)
		cost := (prompt_won + completion_won) * 1.1

		logger.Debug("cost: %f 원", cost)
		logger.Debug("podo: %d 알", int(math.Ceil(cost)))
	}

	return content, nil
}

func (c *Client) RequestChatCompletionsStream(
	ctx context.Context,
	request ChatCompletionRequest,
	onInit func() error,
	onReceiveMessage func(msg string) error,
) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       request.Model,
		"messages":    request.Messages,
		"temperature": request.Temperature,
		"top_p":       request.TopP,
		"stream":      true,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+c.opt.ApiKey)
	req.Header.Set("OpenAI-Organization", c.opt.OrganizationID)
	req.Header.Set("OpenAI-Project", c.opt.ProjectID)

	httpClient := &http.Client{
		Timeout: time.Minute * 5,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	if err := onInit(); err != nil {
		return err
	}

	reader := bufio.NewReader(resp.Body)
	var buffer strings.Builder

	isErrorOnStream := false

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			isErrorOnStream = true
			break
		}

		// stream 응답 스펙
		if len(line) < 6 {
			continue
		}
		if string(line)[:6] != "data: " {
			continue
		}
		line = line[6:]

		if strings.TrimSpace(string(line)) == "[DONE]" {
			if buffer.Len() > 0 {
				msg := buffer.String()
				buffer.Reset()
				onReceiveMessage(msg)
			}
			break
		}

		var responseChunk chatCompletionResponseChunk
		if err := json.Unmarshal(line, &responseChunk); err != nil {
			continue
		}

		for _, choice := range responseChunk.Choices {
			content := choice.Delta.Content
			buffer.WriteString(content)

			if buffer.Len() >= BufferUnit {
				msg := buffer.String()
				buffer.Reset()

				if err := onReceiveMessage(msg); err != nil {
					isErrorOnStream = true
					break
				}
			}
		}
	}

	if buffer.Len() > 0 {
		msg := buffer.String()
		buffer.Reset()

		if err := onReceiveMessage(msg); err != nil {
			isErrorOnStream = true
		}
	}

	if isErrorOnStream {
		return ErrOnStream
	}

	return nil
}

func NewClient() *Client {
	opt := ConstructorOption{
		ApiKey:         config.OpenAiServiceAccountApiKey(),
		OrganizationID: config.OpenAiOrganizationID(),
		ProjectID:      config.OpenAiProjectID(),
	}

	return &Client{
		opt: opt,
	}
}
