package consumer

import (
	"github.com/lwabish/transaction-mapper/pkg/config"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/samber/lo"
	"math"
)

var (
	consumerQianJi = &qianJi{}
)

func init() {
	Registry.Register(consumerQianJi.Name(), func() Plugin {
		return consumerQianJi
	})
}

type qianJi struct{}

func (q qianJi) Name() string {
	return "qianji"
}

func (q qianJi) Transform(transactions []transaction.Transaction, info transaction.AccountInfo) (interface{}, error) {
	result := lo.Map[transaction.Transaction, qianJiTransaction](transactions, func(item transaction.Transaction, index int) qianJiTransaction {
		_, c := config.Config.InferCategory(item)

		var tType, counterpartAccount string
		if toAccount := config.Config.InferTransferToAccount(item, info); toAccount != "" {
			tType = "转账"
			counterpartAccount = toAccount
		} else if item.Amount > 0 {
			tType = "支出"
		} else {
			tType = "收入"
		}

		return qianJiTransaction{
			Time:     item.Time.Format("2006/01/02 15:04"),
			Category: c,
			Type:     tType,
			Amount:   math.Abs(item.Amount),
			Account1: info.Name,
			Account2: counterpartAccount,
			Remark:   item.Description,
		}
	})
	return result, nil
}

type qianJiTransaction struct {
	Time          string  `csv:"时间"`
	Category      string  `csv:"分类"`
	Type          string  `csv:"类型"`
	Amount        float64 `csv:"金额"`
	Account1      string  `csv:"账户1"`
	Account2      string  `csv:"账户2"`
	Remark        string  `csv:"备注"`
	BillMark      string  `csv:"账单标记"`
	ServiceCharge string  `csv:"手续费"`
	Coupon        string  `csv:"优惠券"`
	Tag           string  `csv:"标签"`
	BillImage     string  `csv:"账单图片"`
}
