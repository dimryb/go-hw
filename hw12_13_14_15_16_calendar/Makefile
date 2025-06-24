BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
COMPOSE_FILE := "deployments/docker-compose.yml"
INTEGRATION_COMPOSE_FILE := "deployments/docker-compose.integration.yml"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DOCKER_COMPOSE_TEST := docker compose -f ${INTEGRATION_COMPOSE_FILE}

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) -config ./configs/calendar.yaml

build-service:
	docker build \
        --build-arg LDFLAGS="$(LDFLAGS)" \
        -t $(DOCKER_IMG) \
        -f $(DOCKERFILE_PATH) .

build-calendar-img:
	$(MAKE) build-service \
        DOCKER_IMG=calendar:develop \
        DOCKERFILE_PATH=build/calendar/Dockerfile

build-scheduler-img:
	$(MAKE) build-service \
        DOCKER_IMG=scheduler:develop \
        DOCKERFILE_PATH=build/scheduler/Dockerfile

build-sender-img:
	$(MAKE) build-service \
        DOCKER_IMG=sender:develop \
        DOCKERFILE_PATH=build/sender/Dockerfile

build-testrunner-img:
	docker build \
            -f build/testrunner/Dockerfile \
            -t calendar-testrunner .

up:
	docker compose -f ${COMPOSE_FILE} up -d

down:
	docker compose -f ${COMPOSE_FILE} down

logs:
	docker compose -f ${COMPOSE_FILE} logs -f

rebuild: build-calendar-img build-scheduler-img build-sender-img up

restart: down up

integration-tests:
	docker compose -f ${INTEGRATION_COMPOSE_FILE} down
	docker compose -f ${INTEGRATION_COMPOSE_FILE} up --abort-on-container-exit

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.5

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc \
		-I proto \
		--go_out=proto --go_opt=paths=source_relative \
		--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
		proto/calendar/*.proto

generate-mocks:
	go generate ./...

swagger:
	swag init -g cmd/calendar/main.go -o internal/server/http/docs/

rabbit:
	# http://localhost:15672 guest:guest
	docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

.PHONY: build run build-img build-calendar-img build-scheduler-img build-sender-img run-img version test lint generate generate-mocks swagger rabbit up down logs rebuild restart
