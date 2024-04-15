package bank

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
)

func init() {
	Registry.register("icbc", func() Plugin {
		return &icbc{}
	})
}

type icbc struct {
}

func (i icbc) Name() string {
	return "icbc"
}

func (i icbc) ParseCSV(csvData []byte) ([]transaction.Transaction, error) {
	return []transaction.Transaction{}, nil
}
