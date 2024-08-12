package transaction

import "time"

type Transaction struct {
	Time        time.Time
	Amount      float64
	CNY         bool
	Description string
}

type AccountInfo struct {
	Type string
	Name string
}
