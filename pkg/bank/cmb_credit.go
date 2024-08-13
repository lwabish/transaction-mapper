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
		lo.Map(obj.Data.Detail, func(item DetailBill, index int) transaction.Provider { return item }),
	), nil
}

func (in *cmbCredit) PreProcess(data []byte) (string, error) {
	return string(data), nil
}

type cmbCreditTxn struct {
	RespCode string `json:"respCode"`
	RespMsg  string `json:"respMsg"`
	Data     Data   `json:"data"`
}

type Data struct {
	RMBBillInfo    RMBBillInfo    `json:"rmbBillInfo"`
	DollarBillInfo DollarBillInfo `json:"dollarBillInfo"`
	Detail         []DetailBill   `json:"detail"`
}

type RMBBillInfo struct {
	Amount         string `json:"amount"`
	BillCycleStart string `json:"billCycleStart"`
	BillCycleEnd   string `json:"billCycleEnd"`
	RepaymentDate  string `json:"repaymentDate"`
	MinPayment     string `json:"minPayment"`
	LastBill       string `json:"lastBill"`
	LastRepayment  string `json:"lastRepayment"`
	CreditAmount   string `json:"creditAmount"`
	Interest       string `json:"interest"`
	DebitAmount    string `json:"debitAmount"`
}

type DollarBillInfo struct {
	Amount         string `json:"amount"`
	BillCycleStart string `json:"billCycleStart"`
	BillCycleEnd   string `json:"billCycleEnd"`
	RepaymentDate  string `json:"repaymentDate"`
	MinPayment     string `json:"minPayment"`
	LastBill       string `json:"lastBill"`
	LastRepayment  string `json:"lastRepayment"`
	CreditAmount   string `json:"creditAmount"`
	Interest       string `json:"interest"`
	DebitAmount    string `json:"debitAmount"`
}

type DetailBill struct {
	BillId             string `json:"billId"`
	BillType           string `json:"billType"`
	BillDate           string `json:"billDate"`
	BillMonth          string `json:"billMonth"`
	Org                string `json:"org"`
	TransactionAmount  string `json:"transactionAmount"`
	Amount             string `json:"amount"`
	Description        string `json:"description"`
	PostingDate        string `json:"postingDate"`
	Location           string `json:"location"`
	TotalStages        string `json:"totalStages"`
	CurrentStages      string `json:"currentStages"`
	RemainingStages    string `json:"remainingStages"`
	EffectiveDate      string `json:"effectiveDate"`
	TransactionType    string `json:"transactionType"`
	CardNo             string `json:"cardNo"`
	UniqueNo           string `json:"uniqueNo"`
	CommonDescFlag     string `json:"commonDescFlag"`
	RefundTimeHideFlag string `json:"refundTimeHideFlag"`
	StageOrderNo       string `json:"stageOrderNo"`
	StageTypeFlag      string `json:"stageTypeFlag"`
	SceneLimitType     string `json:"sceneLimitType"`
	SceneLimitAmount   string `json:"sceneLimitAmount"`
}

func (d DetailBill) ParseTime() time.Time {
	t, err := time.Parse("20060102150405", d.BillDate)
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

func (d DetailBill) ParseAmount() float64 {
	return util.ParseFloat(d.TransactionAmount)
}

func (d DetailBill) ParseCNY() bool {
	return true
}

func (d DetailBill) ParseDescription() string {
	return d.Description
}
