package appexport

import (
	"encoding/json"
	"syscall/js"

	"github.com/ZenoN-Cloud/zeno-export/export/xlsx"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) RegisterWASMFunctions() {
	js.Global().Set("zenoExport", js.ValueOf(map[string]interface{}{
		"toExcel": js.FuncOf(a.toExcelWrapper),
		"version": "v0.1.0",
	}))

	// Signal that export WASM is ready
	js.Global().Call("dispatchEvent", js.Global().Get("CustomEvent").New("exportWasmReady"))
}

func (a *App) toExcelWrapper(this js.Value, args []js.Value) interface{} {
	// Return a Promise
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			defer func() {
				if r := recover(); r != nil {
					reject.Invoke(js.ValueOf(map[string]interface{}{
						"error": "internal error during Excel generation",
					}))
				}
			}()

			if len(args) < 3 {
				reject.Invoke(js.ValueOf(map[string]interface{}{
					"error": "missing normalized data argument",
				}))
				return
			}

			// Parse normalized data
			var normalized map[string]interface{}
			if err := json.Unmarshal([]byte(args[2].String()), &normalized); err != nil {
				reject.Invoke(js.ValueOf(map[string]interface{}{
					"error": "failed to parse normalized data: " + err.Error(),
				}))
				return
			}

			// Generate Excel file
			bytes, err := xlsx.BuildFromNormalized(normalized)
			if err != nil {
				reject.Invoke(js.ValueOf(map[string]interface{}{
					"error": "failed to generate Excel: " + err.Error(),
				}))
				return
			}

			// Convert to Uint8Array
			uint8Array := js.Global().Get("Uint8Array").New(len(bytes))
			js.CopyBytesToJS(uint8Array, bytes)

			resolve.Invoke(uint8Array)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
