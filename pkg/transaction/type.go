package transaction

type Transaction struct {
	Date        string
	Amount      float64
	CNY         bool
	Description string
}

type AccountInfo struct {
	Type string
	Name string
}
