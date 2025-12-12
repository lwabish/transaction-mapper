package bank

import (
	"encoding/json"
	"testing"
)

// testCmbCreditData is a minimal sample of CMB credit JSON data for testing
const testCmbCreditData = `{
    "credit_limit": "¥ 50,000.00",
    "payment_due_date": "2025年12月15日",
    "current_balance": "¥ 8,456.78",
    "minimum_payment": "¥ 425.32",
    "statement_date": "2025年11月25日",
    "transaction_details": [
        {
            "card_no": "1234",
            "description": "自动还款",
            "rmb_amount": "-3,256.89"
        },
        {
            "card_no": "5678",
            "description": "滴滴支付-北京某某科技有限公司",
            "original_tran_amount": "-28.45(CN)",
            "posted_date": "11/03",
            "rmb_amount": "-28.45",
            "sold_date": "11/02"
        },
        {
            "card_no": "5678",
            "description": "财付通-滴滴出行",
            "original_tran_amount": "35.67(CN)",
            "posted_date": "11/05",
            "rmb_amount": "35.67",
            "sold_date": "11/04"
        },
        {
            "card_no": "5678",
            "description": "美团支付-深圳某某科技有限公司",
            "original_tran_amount": "12.89(CN)",
            "posted_date": "11/07",
            "rmb_amount": "12.89",
            "sold_date": "11/06"
        }
    ],
    "current_balance_summary": "¥ 8,456.78",
    "balance_b_f": "¥ 5,234.12",
    "payment": "¥ 3,256.89",
    "new_charges": "¥ 8,456.78",
    "adjustment": "¥ 23.45",
    "interest": "¥ 0.00"
}`

func TestCmbCreditPreProcess(t *testing.T) {
	// Use embedded test data
	data := []byte(testCmbCreditData)

	plugin := &cmbCredit{}
	result, err := plugin.PreProcess(data)
	if err != nil {
		t.Fatalf("PreProcess failed: %v", err)
	}

	// Check if result is valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Check if transaction_details exists and doesn't contain 自动还款
	transactionDetails, exists := parsed["transaction_details"]
	if !exists {
		t.Fatal("transaction_details field missing")
	}

	details, ok := transactionDetails.([]interface{})
	if !ok {
		t.Fatal("transaction_details is not an array")
	}

	// Check that no transaction has description = 自动还款
	for i, detail := range details {
		detailMap, ok := detail.(map[string]interface{})
		if !ok {
			t.Fatalf("Transaction detail %d is not an object", i)
		}

		description, ok := detailMap["description"].(string)
		if !ok {
			t.Fatalf("Transaction detail %d missing description", i)
		}

		if description == "自动还款" {
			t.Errorf("Transaction detail %d still contains 自动还款", i)
		}
	}

	// Count original vs filtered transactions
	var originalParsed map[string]interface{}
	err = json.Unmarshal([]byte(testCmbCreditData), &originalParsed)
	if err != nil {
		t.Fatalf("Original data is not valid JSON: %v", err)
	}

	originalDetails := originalParsed["transaction_details"].([]interface{})
	filteredCount := len(details)
	originalCount := len(originalDetails)

	// Should have filtered out exactly 1 "自动还款" transaction
	expectedCount := originalCount - 1
	if filteredCount != expectedCount {
		t.Errorf("Expected %d transactions after filtering, got %d", expectedCount, filteredCount)
	}

	t.Logf("Successfully filtered out 自动还款 transactions: %d -> %d", originalCount, filteredCount)

	// Verify the JSON structure is preserved
	for _, field := range []string{"credit_limit", "payment_due_date", "current_balance", "minimum_payment", "statement_date", "current_balance_summary", "balance_b_f", "payment", "new_charges", "adjustment", "interest"} {
		if _, exists := parsed[field]; !exists {
			t.Errorf("Field %s is missing from processed JSON", field)
		}
	}

	t.Log("All required fields preserved in processed JSON")
}

func TestCmbCreditParse(t *testing.T) {
	// Use embedded test data
	data := []byte(testCmbCreditData)

	plugin := &cmbCredit{}

	// First pre-process to remove 自动还款
	processedData, err := plugin.PreProcess(data)
	if err != nil {
		t.Fatalf("PreProcess failed: %v", err)
	}

	// Then parse the processed data
	transactions, err := plugin.Parse(processedData)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check that we got transactions
	if len(transactions) == 0 {
		t.Fatal("No transactions parsed")
	}

	// Verify no transaction has description = 自动还款
	for _, txn := range transactions {
		if txn.Description == "自动还款" {
			t.Errorf("Transaction with description '自动还款' should have been filtered out")
		}
	}

	t.Logf("Successfully parsed %d transactions", len(transactions))
}

func TestCmbCreditPreProcessOutputValidity(t *testing.T) {
	// Use embedded test data
	data := []byte(testCmbCreditData)

	plugin := &cmbCredit{}
	result, err := plugin.PreProcess(data)
	if err != nil {
		t.Fatalf("PreProcess failed: %v", err)
	}

	// Try to unmarshal and re-marshal to ensure JSON validity
	var parsed interface{}
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Re-marshal to compact format
	compact, err := json.Marshal(parsed)
	if err != nil {
		t.Fatalf("Failed to re-marshal JSON: %v", err)
	}

	// Verify compact version is also valid
	var reParsed interface{}
	err = json.Unmarshal(compact, &reParsed)
	if err != nil {
		t.Fatalf("Compact JSON is not valid: %v", err)
	}

	t.Logf("Generated JSON is valid (size: %d bytes -> %d bytes compact)", len(result), len(compact))
}

func TestCmbCreditTransactionDetailParsing(t *testing.T) {
	// Test individual TransactionDetail parsing functionality
	detail := TransactionDetail{
		SoldDate:           "11/04",
		PostedDate:         "11/05",
		Description:        "测试交易",
		RmbAmount:          "35.67",
		CardNo:             "5678",
		OriginalTranAmount: "35.67(CN)",
	}

	// Test time parsing
	parsedTime := detail.ParseTime()
	if parsedTime.Day() != 4 || parsedTime.Month() != 11 {
		t.Errorf("Expected day 4, month 11, got day %d, month %d", parsedTime.Day(), parsedTime.Month())
	}

	// Test amount parsing
	amount := detail.ParseAmount()
	if amount != 35.67 {
		t.Errorf("Expected amount 35.67, got %f", amount)
	}

	// Test CNY parsing
	if !detail.ParseCNY() {
		t.Error("Expected ParseCNY() to return true")
	}

	// Test description parsing
	desc := detail.ParseDescription()
	if desc != "测试交易" {
		t.Errorf("Expected description '测试交易', got '%s'", desc)
	}

	t.Log("TransactionDetail parsing functionality verified")
}