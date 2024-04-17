package bank

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"os"
)

type Plugin interface {
	Name() string
	PreProcess(csvData []byte) (*os.File, error)
	ParseCSV(*os.File) ([]transaction.Transaction, error)
}
