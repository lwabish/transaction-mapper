package main

import (
	"flag"
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"log"
)

var (
	bankName     string
	consumerName string
)

func init() {
	flag.StringVar(&bankName, "bank", "", "bank identifier")
	flag.StringVar(&consumerName, "consumer", "", "consumer identifier")
}

func main() {
	flag.Parse()
	log.Println("bank:", bankName, "consumer:", consumerName)
	_, err := bank.Registry.Get(bankName)
	if err != nil {
		log.Fatalln(err)
	}
}
