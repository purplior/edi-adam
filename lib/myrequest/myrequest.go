package myrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/purplior/podoroot/domain/shared/logger"
)

type (
	PostRequestOption struct {
		Headers map[string]string
		Body    interface{}
	}
)

func POST(
	url string,
	opt PostRequestOption,
	responseBody interface{},
) (statusCode int, err error) {
	var requestBodyBytes []byte = []byte{}
	if opt.Body != nil {
		requestBodyBytes, err = json.Marshal(opt.Body)
		if err != nil {
			return statusCode, err
		}
	}

	logger.Debug("POST " + url)
	logger.Debug(string(requestBodyBytes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return statusCode, err
	}
	for key, value := range opt.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return statusCode, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	logger.Debug("statusCode: %d", statusCode)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return statusCode, err
	}
	if responseBody != nil {
		if err := json.Unmarshal(body, responseBody); err != nil {
			return statusCode, err
		}
	}

	return statusCode, nil
}
