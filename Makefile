BUILD_DATE := $(shell date +'%d.%m.%Y %H:%M:%S')
BUILD_COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: all
all: clean build

.PHONY: build
build:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/myfoodbot && \
	 go build \
	 -ldflags "-X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	 -o ../../bin/myfoodbot .

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf ./bin		 