FROM golang:1.25.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=js GOARCH=wasm go build -o export.wasm ./cmd/wasm-export
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .

FROM nginx:alpine
COPY --from=builder /app/export.wasm /usr/share/nginx/html/
COPY --from=builder /app/wasm_exec.js /usr/share/nginx/html/
COPY --from=builder /app/demo.html /usr/share/nginx/html/index.html
COPY --from=builder /app/integration-demo.html /usr/share/nginx/html/
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80