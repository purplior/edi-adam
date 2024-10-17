package context

import "context"

type (
	Task                 func(ctx context.Context) (interface{}, error)
	TransactionProcessor func(task Task) (interface{}, error)
)

func SetTransactionProcessor(
	newTransactionProcessor TransactionProcessor,
) {
	transactionProcessor = newTransactionProcessor
}

func Transaction(
	task Task,
) (interface{}, error) {
	if transactionProcessor != nil {
		return transactionProcessor(task)
	}

	return nil, nil
}

var (
	transactionProcessor TransactionProcessor = nil
)
