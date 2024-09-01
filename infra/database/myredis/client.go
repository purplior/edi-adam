package myredis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type (
	Client struct {
		*redis.Client
	}
)

func (c *Client) Connect(ctx context.Context) error {
	_, err := c.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("# [myredis] is connected.")

	return nil
}

func NewClient() *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 서버 주소로 변경
		Password: "",               // 로컬 개발용 비밀번호 없음
		DB:       0,                // 기본 DB
	})

	return &Client{
		client,
	}
}
