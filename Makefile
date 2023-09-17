.PHONY: build
build:
	@mkdir -p build
	@go build -a -ldflags="-X github.com/qase-tms/qasectl/internal.Version=$(shell git describe --tags)" -o build/qli ./main.go

clean:
	@rm -rf ./build/*