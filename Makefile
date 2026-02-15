.PHONY: all build run-server run-tui test clean

BINARY_SERVER=bin/server
BINARY_TUI=bin/tui

all: build

build:
	go build -o $(BINARY_SERVER) ./cmd/server
	go build -o $(BINARY_TUI) ./cmd/tui

run-server:
	go run ./cmd/server

run-tui:
	go run ./cmd/tui

collect:
	curl -X POST http://localhost:8080/api/v1/collect

trends:
	curl http://localhost:8080/api/v1/trends | jq

sources:
	curl http://localhost:8080/api/v1/sources | jq

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf data/profiles/tech/trends/*.md

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

.PHONY: docker
docker:
	docker build -t r3f-signal-agent .
