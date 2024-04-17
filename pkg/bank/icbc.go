package bank

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"log"
	"os"
	"strings"
)

const (
	name = "icbc"
)

func init() {
	Registry.register(name, func() Plugin {
		return &icbc{}
	})
}

type icbc struct {
}

func (i icbc) PreProcess(csvData []byte) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return nil, err
	}

	_, err = tmpFile.WriteString(strings.SplitN(string(csvData), "\n", 7)[6])
	if err != nil {
		return nil, err
	}
	log.Println(tmpFile.Name())
	return tmpFile, nil
}

func (i icbc) Name() string {
	return name
}

func (i icbc) ParseCSV(f *os.File) ([]transaction.Transaction, error) {
	return []transaction.Transaction{}, nil
}
