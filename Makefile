VERSION=0.0.13
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"

all: mackerel-plugin-log-incr-rate

.PHONY: mackerel-plugin-log-incr-rate linux check lint

mackerel-plugin-log-incr-rate: *.go
	go build $(LDFLAGS) -o mackerel-plugin-log-incr-rate

linux: main.go parser.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-log-incr-rate

check:
	go test -v ./...
	go test -race

lint:
	golangci-lint run --timeout 5m ./...