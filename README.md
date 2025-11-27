# zeno-export

GO-WASM Export Engine for Zeno CY  
Local Excel generation from normalized transactions.  
Runs fully in the browser. Zero-retention. No server-side processing.

[![Pipeline](https://gitlab.com/zeno-cy/zeno-export/badges/main/pipeline.svg)](https://gitlab.com/zeno-cy/zeno-export/-/pipelines)
[![Coverage](https://gitlab.com/zeno-cy/zeno-export/badges/main/coverage.svg)](https://gitlab.com/zeno-cy/zeno-export/-/commits/main)

---

## 1. Mission

This module converts normalized statement data (from `zeno-engine`) into downloadable Excel files.

- ✅ 100% local browser execution
- ✅ No data sent to backend
- ✅ Privacy-by-design
- ✅ Independent WASM module
- ✅ Lazy loading only when needed

---

## 2. Architecture

```
[Normalized JSON] → [export.wasm] → [Excel Uint8Array] → [Download]
```

### Components

- `export/xlsx/` - Excel generation engine
- `internal/appexport/` - WASM bridge
- `cmd/wasm-export/` - WASM entry point

---

## 3. Usage

### Build WASM Module

```bash
make build-wasm
```

### JavaScript API

```javascript
// Load export module (lazy)
const excelBytes = await zenoExport.toExcel(JSON.stringify(normalizedData));

// Download Excel file
const blob = new Blob([excelBytes], { 
    type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' 
});
const url = URL.createObjectURL(blob);
const a = document.createElement('a');
a.href = url;
a.download = 'transactions.xlsx';
a.click();
```

### Input Format

```json
{
  "vendor": "hellenic_bank",
  "rows": [
    {
      "booking_date": "2024-11-01",
      "value_date": "2024-11-01",
      "amount": "2800.00",
      "currency": "EUR",
      "description": "INCOMING TRANSFER - SALARY",
      "raw_type": "income",
      "balance": "6200.00"
    }
  ],
  "warnings": []
}
```

---

## 4. Integration with zeno-engine

```javascript
// Step 1: Normalize CSV (core.wasm)
const normalized = await zenoEngine.normalizeCSV(csvData, {});

// Step 2: Generate Excel (export.wasm)
const excelBytes = await zenoExport.toExcel(JSON.stringify(normalized));

// Step 3: Download
downloadFile(excelBytes, 'statement.xlsx');
```

---

## 5. Features

- ✅ Multi-sheet support (Transactions + Warnings)
- ✅ Styled headers and formatting
- ✅ Auto-sized columns
- ✅ Proper number formatting
- ✅ Date formatting
- ✅ Error handling with promises
- ✅ Memory-efficient streaming

---

## 6. Testing

```bash
# Run tests
go test ./...

# Build and serve demo
make serve
# Open http://localhost:8080/demo.html
```

---

## 7. File Structure

```
zeno-export/
├── export/xlsx/          # Excel generation core
│   ├── builder.go        # Main Excel builder
│   ├── types.go          # Data structures
│   └── builder_test.go   # Tests
├── internal/appexport/   # WASM bridge
│   └── wasm.go          # JavaScript interface
├── cmd/wasm-export/     # WASM entry point
│   └── main.go
├── testdata/            # Test fixtures
├── demo.html            # Interactive demo
└── Makefile             # Build automation
```
