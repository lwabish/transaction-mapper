package consumer

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"math"
	"strconv"
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
	var res []blueCoinsTransaction
	for _, t := range transactions {
		bt := blueCoinsTransaction{
			Date:   t.Date,
			Amount: strconv.FormatFloat(math.Abs(t.Amount), 'f', 2, 64),
			Notes:  t.Description,
		}
		// fixme: 转账/退款
		// todo: 分类 货币 账户
		if t.Amount > 0 {
			bt.Type = "i"
		} else {
			bt.Type = "e"
		}

		res = append(res, bt)
	}
	return res, nil
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
