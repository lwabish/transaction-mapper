package engine

import (
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"github.com/lwabish/transaction-mapper/pkg/consumer"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
)

type engine interface {
	ParseBankData(bank.Plugin, []byte) ([]transaction.Transaction, error)

	Transform(consumer.Plugin, []transaction.Transaction) ([]byte, error)
}
