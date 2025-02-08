package consumer

import (
	"fmt"
	"github.com/lwabish/transaction-mapper/pkg/config"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/samber/lo"
	"log"
	"math"
	"strconv"
	"strings"
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

func (b *blueCoins) Name() string {
	return "bluecoins"
}

func (b *blueCoins) Transform(transactions []transaction.Transaction, ai transaction.AccountInfo) (interface{}, error) {
	var res []blueCoinsTransaction
	res = lo.Map[transaction.Transaction, blueCoinsTransaction](transactions, func(item transaction.Transaction, index int) blueCoinsTransaction {
		if !item.CNY {
			log.Printf("found transaction with non-cny currency: %+v", item)
		}
		// 二级分类
		parentCategory, category := config.Config.InferCategory(item)
		bt := blueCoinsTransaction{
			Amount:         strconv.FormatFloat(math.Abs(item.Amount), 'f', 2, 64),
			Notes:          item.Description,
			Date:           item.Time.Format("01-02-2006"),
			ItemOrPayee:    "自动导入",
			ParentCategory: parentCategory,
			Category:       category,
			AccountType:    ai.Type,
			Account:        ai.Name,
		}
		if toAccountType, toAccountName := config.Config.InferTransferToAccount(item, ai); toAccountName != "" {
			bt.Type = "t"
			bt.toAccountType = toAccountType
			bt.toAccountName = toAccountName
			bt.Amount = fmt.Sprintf("-%s", bt.Amount)
			bt.ParentCategory, bt.Category = "(Transfer)", "(Transfer)"
		} else if item.Amount > 0 {
			bt.Type = "i"
		} else if item.Amount < 0 {
			bt.Type = "e"
		}
		return bt
	})
	return b.appendTransferTransaction(res), nil
}

func (b *blueCoins) appendTransferTransaction(bts []blueCoinsTransaction) []blueCoinsTransaction {
	var res []blueCoinsTransaction
	copy(bts, res)
	lo.ForEach[blueCoinsTransaction](bts, func(item blueCoinsTransaction, index int) {
		res = append(res, item)
		if item.Type == "t" {
			log.Printf("found transfer transaction: %+v", item)
			item.Amount = strings.ReplaceAll(item.Amount, "-", "")
			item.AccountType, item.Account = item.toAccountType, item.toAccountName
			res = append(res, item)
		}
	})
	return res
}

// https://drive.google.com/file/d/1rpGg48aZLrhIYRK3JK46M8REOaqJf0C6/view
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
	toAccountType  string
	toAccountName  string
}
