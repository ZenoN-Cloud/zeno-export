# Critical Promise Wrapper Fix

## Problem Fixed
Both `zeno-engine` and `zeno-export` had critical bugs in Promise wrapper functions where inner `args` parameter shadowed outer function arguments.

## Changes Made

### zeno-engine/internal/app/wasm.go
```go
// BEFORE (broken)
func (a *App) normalizeCSVWrapper(this js.Value, args []js.Value) interface{} {
    handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        // args here is [resolve, reject], not CSV data!
        csvData := args[2].String() // WRONG - accessing resolve/reject
    })
}

// AFTER (fixed)
func (a *App) normalizeCSVWrapper(this js.Value, args []js.Value) interface{} {
    // Capture outer arguments first
    csvData := args[0].String()
    optsJSON := ""
    if len(args) > 1 && !args[1].IsNull() {
        optsJSON = args[1].String()
    }
    
    handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
        resolve := promiseArgs[0]
        reject := promiseArgs[1]
        // Use captured csvData and optsJSON
    })
}
```

### zeno-export/internal/appexport/wasm.go
```go
// BEFORE (broken)
func (a *App) toExcelWrapper(this js.Value, args []js.Value) interface{} {
    handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        // args here is [resolve, reject], not JSON data!
        normalizedJSON := args[2].String() // WRONG
    })
}

// AFTER (fixed)
func (a *App) toExcelWrapper(this js.Value, args []js.Value) interface{} {
    // Capture outer argument first
    normalizedJSON := args[0].String()
    
    handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
        resolve := promiseArgs[0]
        reject := promiseArgs[1]
        // Use captured normalizedJSON
    })
}
```

## API Usage (Correct)
```javascript
// zeno-engine
const result = await zenoEngine.normalizeCSV(csvData, optionsJSON);

// zeno-export  
const excelBytes = await zenoExport.toExcel(JSON.stringify(normalizedData));
```

## Status
✅ Both modules rebuilt with fixes
✅ Integration demo works correctly
✅ Data contract between modules is perfect