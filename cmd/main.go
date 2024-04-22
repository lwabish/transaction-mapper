package main

import (
	"flag"
	"github.com/gocarina/gocsv"
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"github.com/lwabish/transaction-mapper/pkg/consumer"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"log"
	"os"
)

var (
	bankName       string
	consumerName   string
	inputFileName  string
	outputFileName string
	accountType    string
	accountName    string
)

func init() {
	flag.StringVar(&bankName, "bank", "", "bank identifier")
	flag.StringVar(&consumerName, "consumer", "", "consumer identifier")
	flag.StringVar(&inputFileName, "input", "templates/icbc.csv", "input file name")
	flag.StringVar(&outputFileName, "output", "output.csv", "output file name")
	flag.StringVar(&accountType, "at", "", "account type")
	flag.StringVar(&accountName, "an", "", "account name")
}

func main() {
	flag.Parse()
	if accountType == "" || accountName == "" {
		flag.Usage()
		log.Fatalln("account type and account name are required")
	}
	log.Println("bank:", bankName, "consumer:", consumerName)

	bankPlugin, err := bank.Registry.Get(bankName)
	if err != nil {
		log.Fatalln(err)
	}

	content, err := os.ReadFile(inputFileName)
	if err != nil {
		log.Fatalln(err)
	}

	csvString, err := bankPlugin.PreProcess(content)
	if err != nil {
		log.Fatalln(err)
	}

	transactions, err := bankPlugin.Parse(csvString)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("example transaction:%+v\n", transactions[0])

	consumerPlugin, err := consumer.Registry.Get(consumerName)
	if err != nil {
		log.Fatalln(err)
	}

	results, err := consumerPlugin.Transform(transactions, transaction.AccountInfo{
		Type: accountType,
		Name: accountName,
	})
	if err != nil {
		log.Fatalln(err)
	}

	outFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	err = gocsv.MarshalFile(results, outFile)
	if err != nil {
		log.Fatalln(err)
	}
}
