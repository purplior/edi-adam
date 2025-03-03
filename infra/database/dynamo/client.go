package dynamo

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/shared/constant"
)

type (
	ConstructorOption struct {
		Phase    constant.Phase
		Region   string // AWS 리전 (예: "us-west-2")
		Endpoint string // 로컬 테스트용 엔드포인트 (필요시)
	}

	Client struct {
		*dynamodb.Client
		opt ConstructorOption
	}
)

func (c *Client) ConnectDB() error {
	isDebugMode := config.DebugMode()

	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(c.opt.Region),
	)
	if err != nil {
		return err
	}

	c.Client = dynamodb.NewFromConfig(cfg)

	if isDebugMode {
		log.Println("[dynamo] 디버그모드가 활성화 되었어요")
	}

	// 연결 테스트 (ListTables 호출)
	if err := c.PingDB(); err != nil {
		return err
	}

	log.Println("[dynamo] 데이터베이스 연결에 성공 했어요")
	return nil
}

func (c *Client) ReconnectDB(maxAttempts int, baseDelay time.Duration) {
	attempts := 0
	var err error

	for {
		err = c.PingDB()
		if err == nil {
			log.Println("[dynamo] 데이터베이스 재연결에 성공 했어요")
			return
		}

		attempts++
		if attempts > maxAttempts {
			log.Fatalf("[dynamo] 데이터베이스 재시도 연결 횟수가 최대를 초과 했어요: %v", err)
			return
		}

		delay := baseDelay * time.Duration(1<<attempts)
		log.Printf("[dynamo] 데이터베이스 연결에 실패 했어요, 재시도 중 %v... (attempt %d/%d)", delay, attempts, maxAttempts)
		time.Sleep(delay)

		if err := c.ConnectDB(); err != nil {
			log.Println("[dynamo] 데이터베이스 연결 재시도 실패:", err)
		}
	}
}

func (c *Client) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Client.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})
	return err
}

func (c *Client) Dispose() error {
	c.Client = nil
	return nil
}

func NewClient() *Client {
	opt := ConstructorOption{
		Phase:    config.Phase(),
		Region:   config.AWSBaseRegion(),
		Endpoint: config.AWSDynamoDBEndpoint(),
	}

	return &Client{
		opt: opt,
	}
}
