.PHONY: build
build:
	@mkdir -p build
	@go build -a -ldflags="-X github.com/qase-tms/qasectl/internal.Version=$(shell git describe --tags)" -o build/qasectl ./main.go

clean:
	@rm -rf ./build/*