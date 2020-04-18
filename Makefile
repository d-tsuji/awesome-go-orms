BIN := awesome-go-orms
BUILD_LDFLAGS := "-s -w"

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
	go vet ./...

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean