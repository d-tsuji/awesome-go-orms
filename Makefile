BIN := awesome-go-orms
BUILD_LDFLAGS := "-s -w"
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE=on

.PHONY: all
all: clean build

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN)

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: lint
lint:
	go get golang.org/x/lint/golint
	go vet ./...
	$(GOBIN)/golint -set_exit_status ./...

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean