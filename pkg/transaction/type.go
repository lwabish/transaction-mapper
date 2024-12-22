package transaction

import (
	"github.com/samber/lo"
	"time"
)

// Provider Provide data to initialize a transaction
type Provider interface {
	ParseTime() time.Time
	ParseAmount() float64
	ParseCNY() bool
	ParseDescription() string
}

// Transaction is an intermediate structure for all banks and all accounting software
type Transaction struct {
	Time time.Time
	// >0 支出; <0 收入
	Amount float64
	CNY    bool
	// todo: 所有的bank plugin实现parsing, consumer即可支持转账
	// 转账涉及到的对方账户
	TransferAccount string
	Description     string
}

func NewTransactionFromProvider(p []Provider) []Transaction {
	return lo.Map(p, func(item Provider, index int) Transaction {
		return Transaction{
			Time:        item.ParseTime(),
			Amount:      item.ParseAmount(),
			CNY:         item.ParseCNY(),
			Description: item.ParseDescription(),
		}
	})
}

type AccountInfo struct {
	Type string
	Name string
}
