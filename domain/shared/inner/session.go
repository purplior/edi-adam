package inner

import (
	"context"

	"gorm.io/gorm"
)

const ()

type (
	// @MODEL {도메인 레벨 컨텍스트}
	Session interface {
		// 세션에 관련된 모든 메모리 해제
		Clear()

		// 인증된 사용자 정보
		Identity() *Identity

		// 사용자 정보 설정
		SetIdentity(identity *Identity)

		// 컨텍스트 객체
		Context() context.Context

		// 트랜잭션 가져오기
		Transaction() *gorm.DB

		// 트랜잭션 시작
		BeginTransaction() error

		// 트랜잭션 커밋
		CommitTransaction() error

		// 트랜잭션 롤백
		RollbackTransaction()

		// 남아있는 트랜잭션이 있는 경우 해제하기
		ReleaseTransaction()

		// 멤버만 허용할 경우 검사
		GuardMemberAuth() error

		// 어드민만 허용할 경우 검사
		GuardAdminAuth() error
	}
)
