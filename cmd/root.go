package cmd

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"github.com/lwabish/transaction-mapper/pkg/consumer"
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootArg = struct {
		input       string
		consumer    string
		bank        string
		account     string
		accountType string
	}{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "transaction-mapper",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		inputFileName, consumerName, bankName := rootArg.input, rootArg.consumer, rootArg.bank
		ai := transaction.AccountInfo{
			Name: rootArg.account,
			Type: rootArg.accountType,
		}
		log.Printf("args: %v", rootArg)

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

		outFile, err := os.OpenFile(fmt.Sprintf("%s-%s.csv", strings.Split(path.Base(inputFileName), ".")[0], consumerName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		err = gocsv.MarshalFile(results, outFile)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.transaction-mapper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&rootArg.input, "input", "i", rootArg.input, "input file path")
	rootCmd.Flags().StringVarP(&rootArg.consumer, "consumer", "c", rootArg.consumer, "consumer app name")
	rootCmd.Flags().StringVarP(&rootArg.bank, "bank", "b", rootArg.bank, "bank name")
	rootCmd.Flags().StringVarP(&rootArg.accountType, "account-type", "z", rootArg.accountType, "[optional]account type")
	rootCmd.Flags().StringVarP(&rootArg.account, "account", "a", rootArg.account, "account name")

	_ = rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagsRequiredTogether("input", "consumer", "bank", "account")
}
