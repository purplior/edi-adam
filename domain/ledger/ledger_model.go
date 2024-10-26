package ledger

import "time"

const (
	LedgerAction_AddBySignUpEvent  LedgerAction = 1
	LedgerAction_ConsumeByAssister LedgerAction = 2
)

type (
	LedgerAction uint

	Ledger struct {
		ID         string       `json:"id"`
		WalletID   string       `json:"walletId"`
		PodoAmount int          `json:"podoAmount"`
		Action     LedgerAction `json:"action"`

		// @see Action에 따라 Reason에 저장되는 값이 다름
		// - 안전하지 않은 방법이지만, 트랜잭션을 보장하는데는 문제 없음. (로그 성격에 더 가까움)
		// - 각 액션별로 디테일한 히스토리 저장이 필요하다면, 새로운 테이블 생성을 추천함.
		// AddBySignUpEvent: ""
		// ConsumeByAssister: ":assisterID"
		Reason string `json:"reason"`

		CreatedAt time.Time `json:"createdAt"`
	}
)
