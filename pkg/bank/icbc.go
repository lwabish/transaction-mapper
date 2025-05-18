package bank

import (
	"fmt"
	"github.com/lwabish/transaction-mapper/pkg/csv"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/lwabish/transaction-mapper/pkg/util"
	"github.com/samber/lo"
	"log"
	"strings"
	"time"
)

var (
	icbcIns = &icbc{}
)

func init() {
	Registry.register(icbcIns.Name(), func() Plugin {
		return icbcIns
	})
}

type icbc struct {
}

func (i *icbc) PreProcess(data []byte) (string, error) {
	withoutHeader := strings.SplitN(string(data), "\n", 7)[6]
	tmp := strings.SplitAfter(withoutHeader, "\n")
	withoutTail := tmp[:len(tmp)-2]
	return strings.ReplaceAll(strings.Join(withoutTail, ""), ",\r\n", "\r\n"), nil
}

func (i *icbc) Name() string {
	return "icbc"
}

func (i *icbc) Parse(data string) ([]transaction.Transaction, error) {

	var ts []icbcTxn
	err := csv.Parse(data, &ts)
	if err != nil {
		log.Fatalln(err)
	}

	return transaction.NewTransactionFromProvider(
		lo.Map(ts, func(item icbcTxn, index int) transaction.Provider { return item }),
	), nil
}

type icbcTxn struct {
	TranDate             string `csv:"交易日期"`
	AccountDate          string `csv:"记账日期"`
	Abstract             string `csv:"摘要"`
	Platform             string `csv:"交易场所"`
	CountryRegion        string `csv:"交易国家或地区简称"`
	TranAmountIncome     string `csv:"交易金额(收入)"`
	TranAmountOutcome    string `csv:"交易金额(支出)"`
	TranCurrency         string `csv:"交易币种"`
	AccountAmountIncome  string `csv:"记账金额(收入)"`
	AccountAmountOutcome string `csv:"记账金额(支出)"`
	AccountCurrency      string `csv:"记账币种"`
	Balance              string `csv:"余额"`
	CounterpartyName     string `csv:"对方户名"`
}

func (i icbcTxn) ParseTime() time.Time {
	t, err := time.Parse("20060102", i.TranDate)
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

func (i icbcTxn) ParseAmount() float64 {
	return lo.Ternary(i.AccountAmountOutcome != "",
		-1*util.ParseFloat(i.AccountAmountOutcome),
		util.ParseFloat(i.AccountAmountIncome),
	)
}

func (i icbcTxn) ParseCNY() bool {
	if i.AccountCurrency == "人民币" {
		return true
	}
	return false
}

func (i icbcTxn) ParseDescription() string {
	return fmt.Sprintf("%s:%s", i.Abstract, i.Platform)
}
