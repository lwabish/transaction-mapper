package consumer

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"log"
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

func (b blueCoins) Transform(transactions []transaction.Transaction, ai transaction.AccountInfo) (interface{}, error) {
	var res []blueCoinsTransaction
	for _, t := range transactions {
		bt := blueCoinsTransaction{
			Date:   t.Date,
			Amount: strconv.FormatFloat(math.Abs(t.Amount), 'f', 2, 64),
			Notes:  t.Description,
		}
		// fixme: 转账/退款
		if t.Amount > 0 {
			bt.Type = "i"
		} else {
			bt.Type = "e"
		}
		// 货币：不填为默认货币，如果非默认货币，需要提供币种和汇率
		if !t.CNY {
			log.Printf("found transaction with non-cny currency: %+v", t)
		}
		// 二级分类，默认记为一日三餐，因为频率最高
		bt.ParentCategory = "食"
		bt.Category = "一日三餐"

		// 二级账户
		bt.AccountType = ai.Type
		bt.Account = ai.Name

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
