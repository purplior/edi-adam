package model

import (
	"time"
)

type (
	WalletLogType uint
)

var (
	WalletLogType_ExpendOnUsingAssister WalletLogType = 1
)

type (
	// @MODEL {지갑의 기록}
	// - 지갑이 진짜 결제되는 금액을 관리하는 지갑이 아니기 때문에
	//   엄격하게 SQL 데이터베이스로 관리하며 데이터의 일관성과 정합성을 보장하지는 않는다.
	// - 최대 1달치만 저장할 예정인데, 몽고 DB의 TTL 인덱스가 유용하다
	// - 대용량 기록일 가능성이 크기 때문에 스케일이 용이한 dynamo DB가 유리하다
	WalletLog struct {
		ID      string        `dynamobav:"id" json:"id"`
		Type    WalletLogType `dynamobav:"type" json:"type"`
		Delta   int           `dynamobav:"delta" json:"delta"`
		Comment string        `dynamobav:"comment" json:"comment"`

		CreatedAt time.Time `dynamobav:"createdAt" json:"createdAt"`

		WalletID uint `dynamobav:"walletId" json:"walletId,omitempty"`
	}
)
