//go:build js && wasm

package main

import (
	"github.com/ZenoN-Cloud/zeno-export/internal/appexport"
)

func main() {
	c := make(chan struct{}, 0)

	// Initialize export app and register WASM functions
	app := appexport.New()
	app.RegisterWASMFunctions()

	<-c
}
