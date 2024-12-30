BIN := "./bin/calendar"
BIN_MIGRATE := "./bin/calendar_migrate"
BIN_GRPC_CLIENT := "./bin/calendar_grpc_client"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/calendar_config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

build-migrate:
	go build -v -o $(BIN_MIGRATE) -ldflags "$(LDFLAGS)" ./cmd/migrate
migrate: build-migrate
	$(BIN_MIGRATE) -config ./configs/calendar_config.toml


build-grpc-client:
	go build -v -o $(BIN_GRPC_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client/grpc

run-grpc-client: build-grpc-client
	$(BIN_GRPC_CLIENT) -config ./configs/calendar_config.toml


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

.PHONY: build run build-img run-img version test lint generate
