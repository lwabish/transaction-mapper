package cmd

import (
	"os"
	"path"
)

type arg struct {
	input       string
	config      string
	App         string `form:"app" binding:"required"`
	Bank        string `form:"bank" binding:"required"`
	Account     string `form:"account" binding:"required"`
	AccountType string `form:"accountType"`
}

func defaultArg() *arg {
	// https://ko.build/features/static-assets/
	return &arg{config: path.Join(os.Getenv("KO_DATA_PATH"), "config.yaml")}
}
