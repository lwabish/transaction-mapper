package bank

import (
	"github.com/lwabish/transaction-mapper/pkg/csv"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"log"
	"strings"
)

const (
	name = "icbc"
)

func init() {
	Registry.register(name, func() Plugin {
		return &icbc{}
	})
}

type icbc struct {
}

func (i icbc) PreProcess(data []byte) (string, error) {
	withoutHeader := strings.SplitN(string(data), "\n", 7)[6]
	tmp := strings.SplitAfter(withoutHeader, "\n")
	withoutTail := tmp[:len(tmp)-3]
	return strings.ReplaceAll(strings.Join(withoutTail, ""), ",\r\n", "\r\n"), nil
}

func (i icbc) Name() string {
	return name
}

func (i icbc) Parse(data string) ([]transaction.Transaction, error) {
	var ts []*t
	err := csv.Parse(data, &ts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ts[0])
	return []transaction.Transaction{}, nil
}

// t 用LLM自动生成结构体 prompt如下
// 我将提供你一个用逗号分隔的csv标题行，请将每一个csv列对应go结构体的一个字段，示例标题行及结构体如下
// 示例标题行：交易日期,记账日期,摘要,交易场所,交易国家或地区简称,交易金额(收入),交易金额(支出)
// 示例结构体：
// 注意：结构体的字段名称请根据标题行翻译成准确的英文，使用首字母大写的驼峰形式，类型全部是string
type t struct {
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
