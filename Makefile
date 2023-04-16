COMMAND_NAME := askGPT
OUTPUT_DIR := ${GOPATH}/bin
OUTPUT_FILE := ${OUTPUT_DIR}/${COMMAND_NAME}
LDFLAGS=-ldflags "-s -w"

.DEFAULT_GOAL: build
.PHONY: build clean vendor test

build: clean vendor
	mkdir -p ${OUTPUT_DIR}
	go build -v -o ${OUTPUT_FILE} ${LDFLAGS}

clean:
	rm -f ${OUTPUT_FILE}

vendor:
	go mod tidy && go mod vendor

test: vendor
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go fmt ./...
	go test ./...
