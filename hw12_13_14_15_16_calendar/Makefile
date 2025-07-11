BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
COMPOSE_FILE := "deployments/docker-compose.yml"
INTEGRATION_COMPOSE_FILE := "deployments/docker-compose.integration.yml"

#GIT_HASH := $(shell git log --format="%h" -n 1)
#LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

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

reset-db:
	docker compose -f ${INTEGRATION_COMPOSE_FILE} down --volumes
	docker compose -f ${INTEGRATION_COMPOSE_FILE} up -d db

integration-tests: reset-db build-testrunner-img
	docker compose -f ${INTEGRATION_COMPOSE_FILE} up --abort-on-container-exit
	docker compose -f ${INTEGRATION_COMPOSE_FILE} down --rmi local

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

.PHONY: build run build-img build-calendar-img build-scheduler-img build-sender-img build-testrunner-img run-img version test lint generate generate-mocks swagger rabbit up down logs rebuild restart integration-tests

# === Настройки Kubernetes ===
CHART_DIR = calendar-chart
RELEASE_NAME = calendar
NAMESPACE = calendar
APP_PORT = 8080
INGRESS_PORT = 80

# === Helm ===

helm-upgrade:
	helm dependency build $(CHART_DIR)
	helm upgrade --install $(RELEASE_NAME) $(CHART_DIR) --namespace $(NAMESPACE) --create-namespace
	kubectl rollout restart deployment -n calendar

helm-install:
	helm install calendar ./calendar-chart --namespace calendar --create-namespace

helm-dependency:
	helm dependency build ./calendar-chart

helm-dependency-update:
	helm dependency update ./calendar-chart

# === Kubernetes ===

k8s-pods:
	kubectl get pods -A

k8s-services:
	kubectl get svc -A

k8s-get-all:
	kubectl get all -A

k8s-logs:
	kubectl logs -f deployment/$(RELEASE_NAME)-api -n $(NAMESPACE)

k8s-logs-scheduler:
	kubectl logs -f deployment/$(RELEASE_NAME)-scheduler -n $(NAMESPACE)

k8s-logs-sender:
	kubectl logs -f deployment/$(RELEASE_NAME)-sender -n $(NAMESPACE)

port-forward:
	kubectl port-forward deployment/calendar-api 8080 -n calendar

k8s-describe-nodes:
	kubectl describe nodes

k8s-check-services:
	kubectl get services -n calendar

.PHONY: helm-upgrade helm-install helm-dependency helm-dependency-update k8s-pods k8s-services k8s-get-all k8s-logs \
k8s-logs-scheduler k8s-logs-sender port-forward k8s-describe-nodes k8s-check-services