.PHONY: run
run: 
	go run ./cmd/main.go

.PHONY: build
build: 
	go build -tags netgo -ldflags '-s -w' -o app ./cmd/main.go

.PHONY: lint
lint:
	golangci-lint run ./...
