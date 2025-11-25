# Compatibility with zeno-engine

## Type Compatibility

The export module now supports both string-based and typed data from zeno-engine:

### Supported Input Types

- **Dates**: `string` (ISO format) or `time.Time`
- **Amounts**: `string`, `float64`, `json.Number`, or `decimal.Decimal`
- **Balance**: `string`, `float64`, `json.Number`, `decimal.Decimal`, or `nil`

### Transaction Validation

All transactions are validated before Excel generation:
- Currency is required
- Description is required  
- Amount must be parseable as float64

### Helper Methods

- `GetBookingDateString()` - converts any date type to string
- `GetValueDateString()` - handles nil and converts to string
- `GetAmountFloat()` - parses amount to float64 with error handling
- `GetBalanceString()` - formats balance with 2 decimal places
- `Validate()` - validates transaction data

## Integration

The integration-demo.html has been updated with proper error handling for WASM loading.

## Testing

Run compatibility tests:
```bash
go test ./export/xlsx/ -v
```