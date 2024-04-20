package consumer

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
)

var (
	consumerBlueCoins = &blueCoins{}
)

func init() {
	Registry.Register(consumerBlueCoins.Name(), func() Plugin {
		return consumerBlueCoins
	})
}

type blueCoins struct {
}

func (b blueCoins) Name() string {
	return "bluecoins"
}

func (b blueCoins) Transform(transactions []transaction.Transaction, path string) error {
	return nil
}
