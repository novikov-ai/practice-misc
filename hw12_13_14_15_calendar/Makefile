BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config_template.toml

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
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

#TODO: implement
#migrate:
	#docker exec -it postgresDB psql -U root
	#shell CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

#TODO: implement
generate:
	rm -rf internal/server/pb
	mkdir -p internal/server/pb

	protoc api/calendar.proto --go_out=internal/server/pb --go-grpc_out=internal/server/pb

	#go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: build run build-img run-img version test lint