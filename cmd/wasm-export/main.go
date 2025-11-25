//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/ZenoN-Cloud/zeno-export/internal/appexport"
)

func main() {
	app := appexport.New()
	app.RegisterWASMFunctions()

	// keep WASM running
	select {}
}
