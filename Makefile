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
lint: $(GOBIN)/golint
	go vet ./...
	$(GOBIN)/golint -set_exit_status ./...

$(GOBIN)/golint:
	cd $(GOBIN) && go get golang.org/x/lint/golint

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean