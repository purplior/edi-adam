package slack

import (
	"github.com/purplior/sbec/application/config"
	"github.com/purplior/sbec/lib/myrequest"
)

var (
	ChannelID_CustomerVoice = "C086YAG8BSR"
)

type (
	SendMessageRequest struct {
		ChannelID string `json:"channel"`
		Text      string `json:"text"`
	}

	ConstructorOption struct {
		APIToken string
	}

	Client struct {
		opt ConstructorOption
	}
)

func (c *Client) SendMessage(request SendMessageRequest) error {
	status, err := myrequest.POST(
		"https://slack.com/api/chat.postMessage",
		myrequest.PostRequestOption{
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + c.opt.APIToken,
			},
			Body: request,
		},
		nil,
	)
	if err != nil || status < 200 || status >= 300 {
		return err
	}

	return nil
}

func NewClient() *Client {
	opt := ConstructorOption{
		APIToken: config.SlackBotToken(),
	}

	return &Client{
		opt: opt,
	}
}
