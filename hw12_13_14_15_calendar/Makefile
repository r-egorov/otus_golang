BINDIR := "./bin"
CALENDAR := "calendar"
MIGRATOR := "migrator"
SCHEDULER := "scheduler"
SENDER := "sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

$(CALENDAR):
	go build -v -o $(BINDIR)/$(CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

$(SCHEDULER):
	go build -v -o $(BINDIR)/$(SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

$(SENDER):
	go build -v -o $(BINDIR)/$(SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

$(MIGRATOR):
	go build -v -o $(BINDIR)/$(MIGRATOR) -ldflags "$(LDFLAGS)" ./cmd/migrator

build: $(CALENDAR) $(SCHEDULER) $(SENDER) $(MIGRATOR)

run: build
	$(CALENDAR) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

up:
	docker-compose --env-file deployments/.env -f deployments/docker-compose.yaml up -d

down:
	docker-compose -f deployments/docker-compose.yaml down

version: $(CALENDAR)
	$(CALENDAR) version

test:
	go test -race ./internal/...

integration-tests:
	docker-compose --env-file deployments/.env.test -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml up -d
	docker logs calendar_tests

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.49.0

lint: install-lint-deps
	golangci-lint run ./...

migrate: $(MIGRATOR)
	$(MIGRATOR) -config=./configs/config.toml -dir=./migrations up

clean:
	rm -rf $(BINDIR)

generate:
	protoc --go_out=./internal/server/ --go-grpc_out=./internal/server/ ./internal/server/grpc/proto/EventService.proto

re: clean build

.PHONY: build run build-img run-img version test lint clean generate re up down
