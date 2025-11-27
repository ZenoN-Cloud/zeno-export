package xlsx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// NormalizedData represents the input from zeno-engine
type NormalizedData struct {
	Vendor   string        `json:"vendor"`
	Rows     []Transaction `json:"rows"`
	Warnings []Warning     `json:"warnings"`
}

// Transaction represents a normalized transaction (compatible with zeno-engine)
type Transaction struct {
	BankVendor  string      `json:"bank_vendor"`
	BookingDate interface{} `json:"booking_date"` // Can be string or time.Time
	ValueDate   interface{} `json:"value_date"`   // Can be string, time.Time, or nil
	Amount      interface{} `json:"amount"`       // Can be string or decimal.Decimal
	Currency    string      `json:"currency"`
	Description string      `json:"description"`
	RawType     string      `json:"raw_type"`
	Balance     interface{} `json:"balance"` // Can be string, decimal.Decimal, or nil
}

// Warning represents a parsing warning
type Warning struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// Validate checks if the transaction has valid data
func (t *Transaction) Validate() error {
	if t.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if t.Description == "" {
		return fmt.Errorf("description is required")
	}
	_, err := t.GetAmountFloat()
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}
	return nil
}

// GetBookingDateString returns booking date as string
func (t *Transaction) GetBookingDateString() string {
	switch v := t.BookingDate.(type) {
	case string:
		return v
	case time.Time:
		return v.Format("2006-01-02")
	default:
		return ""
	}
}

// GetValueDateString returns value date as string
func (t *Transaction) GetValueDateString() string {
	if t.ValueDate == nil {
		return ""
	}
	switch v := t.ValueDate.(type) {
	case string:
		return v
	case time.Time:
		return v.Format("2006-01-02")
	default:
		return ""
	}
}

// GetAmountFloat returns amount as float64
func (t *Transaction) GetAmountFloat() (float64, error) {
	switch v := t.Amount.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	case json.Number:
		return v.Float64()
	default:
		// Handle decimal.Decimal from zeno-engine
		if str := fmt.Sprintf("%v", v); str != "" {
			return strconv.ParseFloat(str, 64)
		}
		return 0, fmt.Errorf("invalid amount type: %T", v)
	}
}

// GetBalanceString returns balance as formatted string
func (t *Transaction) GetBalanceString() string {
	if t.Balance == nil {
		return ""
	}
	switch v := t.Balance.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.2f", v)
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return fmt.Sprintf("%.2f", f)
		}
		return ""
	default:
		// Handle decimal.Decimal from zeno-engine
		if str := fmt.Sprintf("%v", v); str != "" {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return fmt.Sprintf("%.2f", f)
			}
		}
		return ""
	}
}
