BINDIR := "./bin"
BIN := "$(BINDIR)/calendar"
MIGRATOR := "$(BINDIR)/migrator"
SCHEDULER := "$(BINDIR)/scheduler"
SENDER := "$(BINDIR)/sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

$(BIN):
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

$(SCHEDULER):
	go build -v -o $(SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

$(SENDER):
	go build -v -o $(SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

$(MIGRATOR):
	go build -v -o $(MIGRATOR) -ldflags "$(LDFLAGS)" ./cmd/migrator

build: $(BIN) $(SCHEDULER) $(SENDER) $(MIGRATOR)

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: $(BIN)
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate: $(MIGRATOR)
	$(MIGRATOR) -config ./configs/config.toml -dir=./migrations up

clean:
	rm -rf $(BINDIR)

generate:
	protoc --go_out=./internal/server/ --go-grpc_out=./internal/server/ ./internal/server/grpc/proto/EventService.proto

re: clean build

.PHONY: build run build-img run-img version test lint clean generate re
