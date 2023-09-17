build:
	@go build -o bin/gobank ./cmd/api

run: build
	@./bin/gobank

test:
	@go test -v ./...
