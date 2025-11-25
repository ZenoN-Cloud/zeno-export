package xlsx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// BuildFromNormalized converts normalized JSON data into an XLSX byte slice
func BuildFromNormalized(data map[string]interface{}) ([]byte, error) {
	// Parse input data
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var normalized NormalizedData
	if err := json.Unmarshal(jsonBytes, &normalized); err != nil {
		return nil, fmt.Errorf("failed to parse normalized data: %w", err)
	}

	// Create new Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Transactions"
	f.SetSheetName("Sheet1", sheetName)

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E6E6FA"}, Pattern: 1},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %w", err)
	}

	// Set headers
	headers := []string{"Date", "Value Date", "Description", "Amount", "Currency", "Type", "Balance"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 12) // Date
	f.SetColWidth(sheetName, "B", "B", 12) // Value Date
	f.SetColWidth(sheetName, "C", "C", 40) // Description
	f.SetColWidth(sheetName, "D", "D", 15) // Amount
	f.SetColWidth(sheetName, "E", "E", 8)  // Currency
	f.SetColWidth(sheetName, "F", "F", 15) // Type
	f.SetColWidth(sheetName, "G", "G", 15) // Balance

	// Add transaction rows
	for i, tx := range normalized.Rows {
		row := i + 2 // Start from row 2 (after header)

		// Validate transaction
		if err := tx.Validate(); err != nil {
			return nil, fmt.Errorf("invalid transaction at row %d: %w", row, err)
		}

		// Date
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), tx.GetBookingDateString())

		// Value Date
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), tx.GetValueDateString())

		// Description
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), tx.Description)

		// Amount
		amount, err := tx.GetAmountFloat()
		if err != nil {
			return nil, fmt.Errorf("invalid amount in row %d: %w", row, err)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), amount)

		// Currency
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), tx.Currency)

		// Type
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), strings.Title(strings.ReplaceAll(tx.RawType, "_", " ")))

		// Balance
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), tx.GetBalanceString())
	}

	// Add warnings sheet if any
	if len(normalized.Warnings) > 0 {
		warningsSheet := "Warnings"
		f.NewSheet(warningsSheet)

		// Warning headers
		f.SetCellValue(warningsSheet, "A1", "Row")
		f.SetCellValue(warningsSheet, "B1", "Message")
		f.SetCellStyle(warningsSheet, "A1", "B1", headerStyle)

		// Warning data
		for i, warning := range normalized.Warnings {
			row := i + 2
			f.SetCellValue(warningsSheet, fmt.Sprintf("A%d", row), warning.Row)
			f.SetCellValue(warningsSheet, fmt.Sprintf("B%d", row), warning.Message)
		}

		f.SetColWidth(warningsSheet, "A", "A", 8)
		f.SetColWidth(warningsSheet, "B", "B", 50)
	}

	// Generate Excel file as bytes
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write Excel file: %w", err)
	}

	return buffer.Bytes(), nil
}
