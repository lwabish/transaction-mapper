package bank

import (
	"encoding/json"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/lwabish/transaction-mapper/pkg/util"
	"github.com/samber/lo"
	"log"
	"time"
)

type cmbCredit struct{}

var (
	cmbIns = &cmbCredit{}
)

func init() {
	Registry.register(cmbIns.Name(), func() Plugin {
		return cmbIns
	})
}

func (in *cmbCredit) Name() string {
	return "cmbCredit"
}

func (in *cmbCredit) Parse(data string) ([]transaction.Transaction, error) {

	obj := &cmbCreditTxn{}
	if err := json.Unmarshal([]byte(data), obj); err != nil {
		log.Fatalln(err)
	}

	return transaction.NewTransactionFromProvider(
		lo.Map(obj.TransactionDetails, func(item TransactionDetail, index int) transaction.Provider { return item }),
	), nil
}

func (in *cmbCredit) PreProcess(data []byte) (string, error) {
	return string(data), nil
}

type TransactionDetail struct {
	//fixme: 缺少年份信息，目前用当前年份代替
	SoldDate           string `json:"sold_date"`
	PostedDate         string `json:"posted_date"`
	Description        string `json:"description"`
	RmbAmount          string `json:"rmb_amount"`
	CardNo             string `json:"card_no"`
	OriginalTranAmount string `json:"original_tran_amount"`
}

type cmbCreditTxn struct {
	CreditLimit           string              `json:"credit_limit"`
	PaymentDueDate        string              `json:"payment_due_date"`
	CurrentBalance        string              `json:"current_balance"`
	MinimumPayment        string              `json:"minimum_payment"`
	StatementDate         string              `json:"statement_date"`
	TransactionDetails    []TransactionDetail `json:"transaction_details"`
	CurrentBalanceSummary string              `json:"current_balance_summary"`
	BalanceBF             string              `json:"balance_b_f"`
	Payment               string              `json:"payment"`
	NewCharges            string              `json:"new_charges"`
	Adjustment            string              `json:"adjustment"`
	Interest              string              `json:"interest"`
}

func (d TransactionDetail) ParseTime() time.Time {
	currentYear := time.Now().Year()
	parsedTime, err := time.Parse("01/02", d.SoldDate)
	if err != nil {
		log.Fatalln(err)
	}
	yearReplacedTime := time.Date(currentYear, parsedTime.Month(), parsedTime.Day(),
		0, 0, 0, 0, parsedTime.Location())
	return yearReplacedTime
}

func (d TransactionDetail) ParseAmount() float64 {
	return util.ParseFloat(d.RmbAmount)
}

func (d TransactionDetail) ParseCNY() bool {
	return true
}

func (d TransactionDetail) ParseDescription() string {
	return d.Description
}
