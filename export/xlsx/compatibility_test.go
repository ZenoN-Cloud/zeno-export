package xlsx

import (
	"testing"
	"time"
)

func TestCompatibilityWithZenoEngine(t *testing.T) {
	// Simulate data from zeno-engine with time.Time and decimal types
	engineData := map[string]interface{}{
		"vendor": "hellenic_bank",
		"rows": []map[string]interface{}{
			{
				"bank_vendor":  "hellenic_bank",
				"booking_date": "2024-11-01T00:00:00Z", // ISO format from engine
				"value_date":   "2024-11-01T00:00:00Z",
				"amount":       "2800.00", // String representation of decimal
				"currency":     "EUR",
				"description":  "INCOMING TRANSFER - SALARY",
				"raw_type":     "income",
				"balance":      "6200.00",
			},
			{
				"bank_vendor":  "hellenic_bank",
				"booking_date": "2024-11-02T00:00:00Z",
				"value_date":   nil, // nil value date
				"amount":       "-45.50",
				"currency":     "EUR",
				"description":  "CARD PAYMENT - SUPERMARKET LTD",
				"raw_type":     "card_payment",
				"balance":      "6154.50",
			},
		},
		"warnings": []map[string]interface{}{
			{
				"row":     2,
				"message": "Missing value date",
			},
		},
	}

	// Test Excel generation
	excelBytes, err := BuildFromNormalized(engineData)
	if err != nil {
		t.Fatalf("failed to build Excel from engine data: %v", err)
	}

	if len(excelBytes) == 0 {
		t.Error("expected Excel data, got empty bytes")
	}

	t.Logf("Generated Excel from engine data: %d bytes", len(excelBytes))
}

func TestTransactionValidation(t *testing.T) {
	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name: "valid transaction",
			tx: Transaction{
				Currency:    "EUR",
				Description: "Test transaction",
				Amount:      "100.00",
			},
			wantErr: false,
		},
		{
			name: "missing currency",
			tx: Transaction{
				Description: "Test transaction",
				Amount:      "100.00",
			},
			wantErr: true,
		},
		{
			name: "missing description",
			tx: Transaction{
				Currency: "EUR",
				Amount:   "100.00",
			},
			wantErr: true,
		},
		{
			name: "invalid amount",
			tx: Transaction{
				Currency:    "EUR",
				Description: "Test transaction",
				Amount:      "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionHelperMethods(t *testing.T) {
	// Test with time.Time (from zeno-engine)
	now := time.Now()
	tx := Transaction{
		BookingDate: now,
		ValueDate:   now,
		Amount:      "123.45",
		Balance:     "1000.00",
	}

	// Test date methods
	dateStr := tx.GetBookingDateString()
	if dateStr != now.Format("2006-01-02") {
		t.Errorf("GetBookingDateString() = %v, want %v", dateStr, now.Format("2006-01-02"))
	}

	valueDateStr := tx.GetValueDateString()
	if valueDateStr != now.Format("2006-01-02") {
		t.Errorf("GetValueDateString() = %v, want %v", valueDateStr, now.Format("2006-01-02"))
	}

	// Test amount method
	amount, err := tx.GetAmountFloat()
	if err != nil {
		t.Errorf("GetAmountFloat() error = %v", err)
	}
	if amount != 123.45 {
		t.Errorf("GetAmountFloat() = %v, want %v", amount, 123.45)
	}

	// Test balance method
	balance := tx.GetBalanceString()
	if balance != "1000.00" {
		t.Errorf("GetBalanceString() = %v, want %v", balance, "1000.00")
	}

	// Test with nil value date
	tx.ValueDate = nil
	valueDateStr = tx.GetValueDateString()
	if valueDateStr != "" {
		t.Errorf("GetValueDateString() with nil = %v, want empty string", valueDateStr)
	}
}
