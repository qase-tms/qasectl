version := $(shell git describe --tags)

.PHONY: build
build:
	@mkdir -p build
	@go build -a -ldflags="-X github.com/qase-tms/qasectl/internal.Version=$(version)" -o build/qli ./main.go

clean:
	@rm -rf ./build/qli

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

.PHONY: docker
docker:
	@docker build -t ghcr.io/qase-tms/qase:$(version) -t ghcr.io/qase-tms/qase:$latest -f ./build/Dockerfile --build-arg VERSION=$(version) .
