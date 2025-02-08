package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"strings"
)

var (
	rootArg = defaultArg()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "transaction-mapper",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		inputFileName := rootArg.input
		content, err := os.ReadFile(inputFileName)
		if err != nil {
			log.Fatalln(err)
		}
		dst := fmt.Sprintf("%s-%s.csv", strings.Split(path.Base(inputFileName), ".")[0], rootArg.Consumer)
		if err := runEngine(rootArg, content, dst); err != nil {
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
	rootCmd.Flags().StringVarP(&rootArg.Consumer, "consumer", "c", rootArg.Consumer, "consumer app name")
	rootCmd.Flags().StringVarP(&rootArg.Bank, "bank", "b", rootArg.Bank, "bank name")
	rootCmd.Flags().StringVarP(&rootArg.AccountType, "account-type", "z", rootArg.AccountType, "[optional]account type")
	rootCmd.Flags().StringVarP(&rootArg.Account, "account", "a", rootArg.Account, "account name")
	rootCmd.Flags().StringVarP(&rootArg.config, "config", "f", rootArg.config, "config file path")

	_ = rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagsRequiredTogether("input", "consumer", "bank", "account")
}
