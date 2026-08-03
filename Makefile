.PHONY: fmt lint vet test cover bench fuzz wasm

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

test:
	go test -v ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

bench:
	go test -bench=. -benchmem ./...

fuzz:
	go test -fuzz=Fuzz -fuzztime=30s ./...

wasm:
	go build -o dashboard/public/picomatch.wasm cmd/wasm/main.go
