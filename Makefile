build-wasm:
	GOOS=js GOARCH=wasm go build -o export.wasm ./cmd/wasm-export

wasm-exec:
	cp "/opt/homebrew/Cellar/go/1.25.4/libexec/misc/wasm/wasm_exec.js" .

wasm: build-wasm wasm-exec
