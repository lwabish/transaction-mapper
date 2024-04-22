package consumer

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
)

type Plugin interface {
	Name() string
	Transform(transactions []transaction.Transaction, info transaction.AccountInfo) (interface{}, error)
}
