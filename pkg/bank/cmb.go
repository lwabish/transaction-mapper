package bank

import (
	"fmt"
	"github.com/lwabish/transaction-mapper/pkg/csv"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/lwabish/transaction-mapper/pkg/util"
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
	var ts []*cmbTxn
	err := csv.Parse(data, &ts)
	if err != nil {
		log.Fatalln(err)
	}
	var result []transaction.Transaction
	for _, s := range ts {
		tran := transaction.Transaction{}
		t, err := time.Parse("20060102 15:04:05", fmt.Sprintf("%s %s", s.TranDate, s.TranTime))
		if err != nil {
			log.Println(err)
		}
		tran.Time = t
		if s.Outcome != "" {
			tran.Amount = -1 * util.ParseFloat(s.Outcome)
		} else {
			tran.Amount = util.ParseFloat(s.Income)
		}
		tran.CNY = true
		tran.Description = fmt.Sprintf("%s:%s", s.TranType, s.TranRemark)
		result = append(result, tran)
	}
	return result, nil
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
