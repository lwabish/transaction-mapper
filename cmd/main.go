package main

import (
	"flag"
	"github.com/gocarina/gocsv"
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"github.com/lwabish/transaction-mapper/pkg/consumer"
	"log"
	"os"
)

var (
	bankName       string
	consumerName   string
	inputFileName  string
	outputFileName string
)

func init() {
	flag.StringVar(&bankName, "bank", "", "bank identifier")
	flag.StringVar(&consumerName, "consumer", "", "consumer identifier")
	flag.StringVar(&inputFileName, "input", "templates/icbc.csv", "input file name")
	flag.StringVar(&outputFileName, "output", "output.csv", "output file name")
}

func main() {
	flag.Parse()
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

	log.Println(transactions[0])

	consumerPlugin, err := consumer.Registry.Get(consumerName)
	if err != nil {
		log.Fatalln(err)
	}

	results, err := consumerPlugin.Transform(transactions)
	if err != nil {
		log.Fatalln(err)
	}

	err = gocsv.MarshalFile(results, os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
