package main

import (
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"github.com/lwabish/transaction-mapper/pkg/consumer"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"log"
	"os"
	"path"
	"strings"
)

var (
	consumerName  string
	inputFileName string
)

func init() {
	flag.StringVar(&consumerName, "consumer", "", "consumer identifier")
	flag.StringVar(&inputFileName, "input", "", "input file name")
}

func main() {
	flag.Parse()
	if inputFileName == "" {
		flag.Usage()
		log.Fatalln("inputFileName required")
	}
	fileName := path.Base(inputFileName)
	bankName, ai := parseInfoFromFileName(fileName)

	log.Println("bank:", bankName, "consumer:", consumerName)

	bankPlugin, err := bank.Registry.Get(bankName)
	if err != nil {
		log.Fatalln(err)
	}

	content, err := os.ReadFile(inputFileName)
	if err != nil {
		log.Fatalln(err)
	}

	preProcessStr, err := bankPlugin.PreProcess(content)
	if err != nil {
		log.Fatalln(err)
	}

	transactions, err := bankPlugin.Parse(preProcessStr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("example transaction:%+v\n", transactions[0])

	consumerPlugin, err := consumer.Registry.Get(consumerName)
	if err != nil {
		log.Fatalln(err)
	}

	results, err := consumerPlugin.Transform(transactions, ai)
	if err != nil {
		log.Fatalln(err)
	}

	outFile, err := os.OpenFile(fmt.Sprintf("%s-%s.csv", strings.Split(fileName, ".")[0], consumerName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	err = gocsv.MarshalFile(results, outFile)
	if err != nil {
		log.Fatalln(err)
	}
}

// todo: 参数传入
func parseInfoFromFileName(fileName string) (string, transaction.AccountInfo) {
	log.Printf("parsing info from file: %s", fileName)
	parts := strings.SplitN(fileName, "-", 4)
	return parts[0], transaction.AccountInfo{
		Type: parts[1],
		Name: parts[2],
	}
}
