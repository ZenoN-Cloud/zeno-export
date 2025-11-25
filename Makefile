.PHONY: build-wasm clean test

# Build WASM module
build-wasm:
	GOOS=js GOARCH=wasm go build -o export.wasm ./cmd/wasm-export
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" .

# Clean build artifacts
clean:
	rm -f export.wasm wasm_exec.js

# Test the module
test:
	go test ./...

# Build and serve for testing
serve: build-wasm
	python3 -m http.server 8080

# Docker commands
docker-build:
	docker build -t zeno-export .

docker-run: docker-build
	docker run -p 8080:80 zeno-export

docker-compose-up:
	docker-compose up --build

docker-compose-down:
	docker-compose down

# Default target
all: build-wasm