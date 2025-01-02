package podosms

import (
	"fmt"
	"time"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/lib/dt"
	"github.com/purplior/podoroot/lib/mycrypto"
	"github.com/purplior/podoroot/lib/mydate"
	"github.com/purplior/podoroot/lib/myrequest"
)

type (
	SendSMSRequest struct {
		Subject     string
		Content     string
		ToList      []string
		ReserveTime time.Time
	}

	SendSMSResponse struct {
		RequestID   string `json:"requestId"`
		RequestTime string `json:"requestTime"`
		StatusCode  string `json:"statusCode"`
		StatusName  string `json:"statusName"`
	}

	ConstructorOption struct {
		From      string
		ServiceID string
		AccessKey string
		SecretKey string
	}

	Client struct {
		opt ConstructorOption
	}
)

func (c *Client) SendSMS(request SendSMSRequest) (SendSMSResponse, error) {
	contentLen := len(request.Content)
	if contentLen < 4 || contentLen > 40 {
		return SendSMSResponse{}, ErrInvalidContentSize
	}

	url := fmt.Sprintf("https://sens.apigw.ntruss.com/sms/v2/services/%s/messages", c.opt.ServiceID)

	messages := make([]map[string]interface{}, len(request.ToList))
	for i, to := range request.ToList {
		messages[i] = map[string]interface{}{
			"to": to,
		}
	}
	requestBody := map[string]interface{}{
		"type":     "SMS",
		"from":     c.opt.From,
		"content":  request.Content,
		"messages": messages,
	}
	if len(request.Subject) > 0 {
		requestBody["subject"] = request.Subject
	}
	if !request.ReserveTime.IsZero() {
		requestBody["reserveTime"] = request.ReserveTime.Format("2006-01-02 15:04")
	}

	sign, err := mycrypto.SignMapDataWithHMACSHA256(requestBody, c.opt.SecretKey)
	if err != nil {
		return SendSMSResponse{}, ErrSendSMS
	}

	headers := map[string]string{
		"Content-Type":             "application/json",
		"x-ncp-apigw-timestamp":    dt.Str(mydate.NowUnixMilli()),
		"x-ncp-iam-access-key":     c.opt.AccessKey,
		"x-ncp-apigw-signature-v2": sign,
	}

	var responseBody SendSMSResponse
	_, err = myrequest.POST(
		url,
		myrequest.PostRequestOption{
			Headers: headers,
		},
		&responseBody,
	)
	if err != nil {
		return SendSMSResponse{}, ErrSendSMS
	}

	return responseBody, nil
}

func NewClient() *Client {
	opt := ConstructorOption{
		From:      config.NCloudSMSFrom(),
		ServiceID: config.NCloudSMSServiceID(),
		AccessKey: config.NCloudAccessKey(),
		SecretKey: config.NCloudSecretKey(),
	}

	return &Client{
		opt: opt,
	}
}
