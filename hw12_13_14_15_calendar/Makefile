BIN := "./bin/calendar"
BIN_SENDER := "./bin/sender"
BIN_SCHEDULER := "./bin/scheduler"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run_calendar: build_calendar
	$(BIN) -config ./configs/calendar_config.toml

build_sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run_sender: build_sender
	$(BIN_SENDER) -config ./configs/sender_config.toml

build_scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

run_scheduler: build_scheduler
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.toml

build: build_calendar build_scheduler build_sender

run: build
	$(BIN) -config ./configs/calendar_config.toml

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...



docker-build:
	docker compose -f deployments/docker-compose.yml build
docker-env-up:
	docker compose -f deployments/docker-compose.yml up -d cal_database cal_queue
docker-migrate:
	docker compose -f deployments/docker-compose.yml up cal_migrate
up:
	make docker-build
	make docker-env-up
	make docker-migrate
	docker compose -f deployments/docker-compose.yml up -d cal_server cal_scheduler cal_sender
down:
	docker compose -f deployments/docker-compose.yml down


docker-build-test:
	docker compose -f deployments/docker-compose.test.yml build
docker-env-up-test:
	docker compose -f deployments/docker-compose.test.yml up -d cal_database_test cal_queue_test
docker-run-services-test:
	docker compose -f deployments/docker-compose.test.yml up -d cal_server_test cal_scheduler_test cal_sender_test
docker-stop-test:
	docker compose -f deployments/docker-compose.test.yml down

docker-migrate-test:
	docker compose -f deployments/docker-compose.test.yml up cal_migrate_test

integration-test:
	make docker-build-test
	make docker-env-up-test
	make docker-migrate-test
	make docker-run-services-test
	docker compose -f deployments/docker-compose.test.yml up calendar_test
	make docker-stop-test


generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		api/*.proto

	protoc -I ./api \
	--grpc-gateway_out internal/server/grpc/pb \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    api/*.proto

.PHONY: build run build_calendar run_calendar build_sender run_sender build_scheduler run_scheduler build-img run-img version test lint generate
