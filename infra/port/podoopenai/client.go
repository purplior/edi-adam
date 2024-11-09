package podoopenai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/podossaem/podoroot/application/config"
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
		Model    string
		Messages []map[string]string
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
	BufferUnit = 40
)

func (c *Client) RequestChatCompletions(
	ctx context.Context,
	request ChatCompletionRequest,
	onInit func() error,
) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":    request.Model,
		"messages": request.Messages,
	})
	if err != nil {
		return "", err
	}

	fmt.Println(request.Messages)

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
	}

	if err := json.Unmarshal(responseBody, &completionResponse); err != nil {
		return "", err
	}

	if len(completionResponse.Choices) == 0 {
		return "", fmt.Errorf("no completion choices found")
	}

	return completionResponse.Choices[0].Message.Content, nil
}

func (c *Client) RequestChatCompletionsStream(
	ctx context.Context,
	request ChatCompletionRequest,
	onInit func() error,
	onReceiveMessage func(msg string) error,
) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":    request.Model,
		"messages": request.Messages,
		"stream":   true,
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
