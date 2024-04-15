package bank

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
)

type Plugin interface {
	Name() string
	ParseCSV(csvData []byte) ([]transaction.Transaction, error)
}
