package bank

import (
	"fmt"
	"github.com/lwabish/transaction-mapper/pkg/csv"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/lwabish/transaction-mapper/pkg/util"
	"github.com/samber/lo"
	"log"
	"time"
)

type spdb struct{}

func (s spdb) Name() string {
	return "spdb"
}

func (s spdb) PreProcess(data []byte) (string, error) {
	return string(data), nil
}

func (s spdb) Parse(data string) ([]transaction.Transaction, error) {
	var ts []spdbTxn
	err := csv.Parse(data, &ts)
	if err != nil {
		log.Fatalln(err)
	}
	return transaction.NewTransactionFromProvider(lo.Map(ts, func(item spdbTxn, index int) transaction.Provider { return item })), nil
}

var (
	bankSpdb = &spdb{}
)

func init() {
	Registry.register(bankSpdb.Name(), func() Plugin {
		return bankSpdb
	})
}

type spdbTxn struct {
	Date               string `csv:"交易日期Date"`
	Time               string `csv:"交易时间Time"`
	TransactionAccount string `csv:"交易账号TransactionAccount"`
	TransactionName    string `csv:"交易名称TransactionName"`
	TransactionAmount  string `csv:"交易金额TransactionAmount"`
	Balance            string `csv:"账户余额Balance"`
	CounterParty       string `csv:"对手姓名CounterParty"`
	OpponentAccount    string `csv:"对手账号OpponentAccount"`
	Summary            string `csv:"交易摘要Summary"`
}

func (s spdbTxn) ParseTime() time.Time {
	t, err := time.Parse("20060102", s.Date)
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

func (s spdbTxn) ParseAmount() float64 {
	return -1 * util.ParseFloat(s.TransactionAmount)
}

func (s spdbTxn) ParseCNY() bool {
	return true
}

func (s spdbTxn) ParseDescription() string {
	return fmt.Sprintf("%s-%s(%s)", s.CounterParty, s.TransactionName, s.Summary)
}
