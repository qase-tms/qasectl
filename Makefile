.PHONY: build
build:
	@mkdir -p build
	@go build -a -ldflags="-X github.com/qase-tms/qasectl/internal.Version=$(shell git describe --tags)" -o build/qli ./main.go

clean:
	@rm -rf ./build/*

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: generate
generate:
	@go generate ./...

.PHONY: test
test:
	@go test -v ./...

.PHONY: install
install:
	@go mod tidy
	@go install go.uber.org/mock/mockgen@latest

.PHONY: coverage
coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
