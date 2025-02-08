package cmd

import (
	"github.com/gocarina/gocsv"
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"github.com/lwabish/transaction-mapper/pkg/config"
	"github.com/lwabish/transaction-mapper/pkg/consumer"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"log"
	"os"
)

func runEngine(a *arg, inputFile []byte, dst string) error {
	consumerName, bankName := a.Consumer, a.Bank
	ai := transaction.AccountInfo{
		Name: a.Account,
		Type: a.AccountType,
	}
	log.Printf("args: %v", a)

	log.Printf("loading config from file: %s", a.config)
	bs, err := os.ReadFile(a.config)
	if err != nil {
		return err
	}
	conf, err := config.NewConfig(bs)
	if err != nil {
		return err
	}

	bankPlugin, err := bank.Registry.Get(bankName)
	if err != nil {
		return err
	}

	preProcessStr, err := bankPlugin.PreProcess(inputFile)
	if err != nil {
		return err
	}

	transactions, err := bankPlugin.Parse(preProcessStr)
	if err != nil {
		return err
	}

	log.Printf("example transaction:%+v\n", transactions[0])

	consumerPlugin, err := consumer.Registry.Get(consumerName, conf)
	if err != nil {
		return err
	}

	results, err := consumerPlugin.Transform(transactions, ai)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	err = gocsv.MarshalFile(results, outFile)
	if err != nil {
		return err
	}
	return nil
}
