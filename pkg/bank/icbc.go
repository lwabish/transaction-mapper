package bank

import (
	"github.com/gocarina/gocsv"
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
	err := gocsv.UnmarshalString(data, &ts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ts[0])
	return []transaction.Transaction{}, nil
}

type t struct {
	D string `csv:"交易日期"`
}
