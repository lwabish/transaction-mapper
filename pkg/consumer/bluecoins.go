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

func (b blueCoins) Transform(transactions []transaction.Transaction) (interface{}, error) {
	return []blueCoinsTransaction{}, nil
}

type blueCoinsTransaction struct {
	Type           string `csv:"(1)Type"`
	Date           string `csv:"(2)Date"`
	ItemOrPayee    string `csv:"(3)Item or Payee"`
	Amount         string `csv:"(4)Amount"`
	Currency       string `csv:"(5)Currency"`
	ConversionRate string `csv:"(6)ConversionRate"`
	ParentCategory string `csv:"(7)Parent Category"`
	Category       string `csv:"(8)Category"`
	AccountType    string `csv:"(9)Account Type"`
	Account        string `csv:"(10)Account"`
	Notes          string `csv:"(11)Notes"`
	Label          string `csv:"(12) Label"`
	Status         string `csv:"(13) Status"`
	Split          string `csv:"(14) Split"`
}
