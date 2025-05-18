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
	var ts []cmbTxn
	err := csv.Parse(data, &ts)
	if err != nil {
		log.Fatalln(err)
	}
	return transaction.NewTransactionFromProvider(lo.Map(ts, func(item cmbTxn, index int) transaction.Provider { return item })), nil
}

func (in *cmb) PreProcess(data []byte) (string, error) {
	preSplit := strings.SplitN(string(data), "\n", 8)
	withoutHeader := preSplit[7]
	postSplit := strings.SplitAfter(withoutHeader, "\r\n")
	withoutTail := postSplit[:len(postSplit)-4]
	return strings.Join(withoutTail, ""), nil
}

type cmbTxn struct {
	TranDate   string `csv:"交易日期"`
	TranTime   string `csv:"交易时间"`
	Income     string `csv:"收入"`
	Outcome    string `csv:"支出"`
	Balance    string `csv:"余额"`
	TranType   string `csv:"交易类型"`
	TranRemark string `csv:"交易备注"`
}

func (c cmbTxn) ParseTime() time.Time {
	t, err := time.Parse("20060102 15:04:05", fmt.Sprintf("%s %s", c.TranDate, c.TranTime))
	if err != nil {
		log.Println(err)
	}
	return t
}

func (c cmbTxn) ParseAmount() float64 {
	if c.Outcome != "" {
		return util.ParseFloat(c.Outcome)
	} else {
		return -1 * util.ParseFloat(c.Income)
	}
}

func (c cmbTxn) ParseCNY() bool {
	return true
}

func (c cmbTxn) ParseDescription() string {
	return fmt.Sprintf("%s:%s", c.TranType, c.TranRemark)
}
