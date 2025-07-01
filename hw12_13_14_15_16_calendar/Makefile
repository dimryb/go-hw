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
MINIKUBE_PROFILE := calendar-cluster

check-docker:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "Docker не установлен"; exit 1; }

k8s-start: check-docker
	minikube start --profile $(MINIKUBE_PROFILE) --driver=docker
	minikube addons enable ingress --profile $(MINIKUBE_PROFILE)

k8s-stop:
	minikube stop --profile $(MINIKUBE_PROFILE)

k8s-delete:
	minikube delete --profile $(MINIKUBE_PROFILE)

k8s-status:
	minikube status --profile $(MINIKUBE_PROFILE)

k8s-load:
	minikube image load calendar:develop --profile $(MINIKUBE_PROFILE)
	minikube image load scheduler:develop --profile $(MINIKUBE_PROFILE)
	minikube image load sender:develop --profile $(MINIKUBE_PROFILE)

k8s-build-and-load: build-calendar-img build-scheduler-img build-sender-img
	make k8s-load

k8s-check-images:
	minikube ssh --profile $(MINIKUBE_PROFILE) "docker images"

k8s-describe-nodes:
	kubectl describe nodes

k8s-get-all:
	kubectl get all -A

k8s-pods:
	kubectl get pods -A

k8s-apply-deployments:
	kubectl apply -f k8s/calendar/
	kubectl apply -f k8s/scheduler/
	kubectl apply -f k8s/sender/

k8s-apply-config:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secret.yaml

helm-install-postgres:
	helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql \
          --namespace calendar \
          --values k8s/postgres/values.yaml

helm-delete-postgres:
	helm delete postgres -n calendar

k8s-pg-connect:
	$(eval PGPASSWORD := $(shell kubectl get secret -n calendar calendar-secrets -o jsonpath='{.data.DB_PASSWORD}' | base64 --decode))
	kubectl run pg-client --rm -i --tty --image=postgres \
        --namespace calendar --restart=Never \
        --env="PGPASSWORD=$(PGPASSWORD)" \
        --command -- psql -h postgres-postgresql -U user -d calendar

k8s-restart:
	kubectl rollout restart deployment/calendar-api -n calendar

deploy: k8s-load k8s-apply-config helm-install-postgres
	make k8s-apply-deployments

.PHONY: check-docker k8s-start k8s-stop k8s-delete k8s-status k8s-load k8s-build-and-load k8s-check-images \
k8s-describe-nodes k8s-get-all k8s-pods k8s-apply-deployments k8s-apply-config k8s-instal k8s-pg-connect k8s-restart \
helm-install-postgres helm-delete-postgres deploy