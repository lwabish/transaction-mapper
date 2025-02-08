package cmd

type arg struct {
	input       string
	config      string
	App         string `form:"app" binding:"required"`
	Bank        string `form:"bank" binding:"required"`
	Account     string `form:"account" binding:"required"`
	AccountType string `form:"accountType"`
}

func defaultArg() *arg {
	return &arg{config: "config.yaml"}
}
