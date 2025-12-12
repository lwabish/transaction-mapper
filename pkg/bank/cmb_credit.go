package bank

import (
	"encoding/json"
	"log"
	"time"

	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/lwabish/transaction-mapper/pkg/util"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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
	result := string(data)

	// 获取 transaction_details 数组
	transactionDetails := gjson.Get(result, "transaction_details")
	if !transactionDetails.Exists() {
		return result, nil
	}

	// 过滤掉 description=自动还款 的交易
	var filteredDetails []interface{}
	transactionDetails.ForEach(func(key, value gjson.Result) bool {
		description := value.Get("description").String()
		if description != "自动还款" {
			// 将 gjson.Result 转换为 map[string]interface{}
			filteredDetails = append(filteredDetails, value.Value())
		}
		return true
	})

	// 将过滤后的数组重新设置回 JSON
	var err error
	result, err = sjson.Set(result, "transaction_details", filteredDetails)
	if err != nil {
		log.Printf("Failed to update transaction_details: %v", err)
		return string(data), nil
	}

	return result, nil
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
	if d.SoldDate == "" {
		d.SoldDate = d.PostedDate
	}
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
