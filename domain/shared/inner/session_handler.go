package inner

type (
	// @MODULE {도메인 레벨 컨}
	// 세션을 내부 상태로 관리하지는 않음
	// 따라서 Manager가 아닌 Handler
	SessionHandler interface {
		// 트랜잭션 시작
		BeginTransaction(
			session Session,
		) error

		// 트랜잭션 커밋
		CommitTransaction(
			session Session,
		) error

		// 트랜잭션 롤백
		RollbackTransaction(
			session Session,
		)

		// 트랜잭션 초기화
		ReleaseTransaction(
			session Session,
		)
	}
)
