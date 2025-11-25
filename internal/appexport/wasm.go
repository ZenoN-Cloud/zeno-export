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
}

func (a *App) toExcelWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf("Error: missing argument")
	}

	var normalized map[string]interface{}
	if err := json.Unmarshal([]byte(args[0].String()), &normalized); err != nil {
		return js.ValueOf("Error parsing JSON: " + err.Error())
	}

	bytes, err := xlsx.BuildFromNormalized(normalized)
	if err != nil {
		return js.ValueOf("Error building Excel: " + err.Error())
	}

	uint8Array := js.TypedArrayOf(bytes)
	defer uint8Array.Release()
	return uint8Array
}
