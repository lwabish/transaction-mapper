package bank

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
)

type Plugin interface {
	Name() string
	PreProcess(data []byte) (string, error)
	Parse(data string) ([]transaction.Transaction, error)
}
