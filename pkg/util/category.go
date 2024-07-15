package util

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"strings"
)

var (
	categoryKeywords = map[string][]string{
		"行|打车开车": {"打车", "出行"},
		"住|房租水电": {"水费", "电费", "燃气费"},
		"行|物流快递": {"菜鸟", "顺丰"},
	}
	keywordsCategory = map[string]string{}
)

func init() {
	for key, value := range categoryKeywords {
		for _, v := range value {
			keywordsCategory[v] = key
		}
	}
}

func InferCategory(t transaction.Transaction) (string, string) {
	parentCategory, category := inferByRules(t.Description)
	return parentCategory, category
}

func inferByRules(desc string) (string, string) {
	for k, v := range keywordsCategory {
		if strings.Contains(desc, k) {
			parts := strings.Split(v, "|")
			return parts[0], parts[1]
		}
	}
	return "食", "一日三餐"
}
