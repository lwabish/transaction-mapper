package bank

import "github.com/lwabish/transaction-mapper/pkg/transaction"

type cmb struct {
}

var (
	bankCmb = &cmb{}
)

func init() {
	Registry.register(bankCmb.Name(), func() Plugin {
		return bankCmb
	})
}

func (in *cmb) Name() string {
	return "cmb"
}

func (in *cmb) Parse(data string) ([]transaction.Transaction, error) {
	panic("implement me")
}

func (in *cmb) PreProcess(data []byte) (string, error) {
	return string(data), nil
}
