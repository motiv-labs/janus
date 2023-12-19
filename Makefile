NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

VERSION ?= "dev-$(shell git rev-parse --short HEAD)"
GO_LINKER_FLAGS=-ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: all lint test-unit test-integration test-features build

all: test-unit build

# this make target is for testing purposes, run it and then run lms, it will be using local built janus
build-lms:
	docker build -t janus-gateway .
	docker tag janus-gateway krelms/janus-gateway:latest

# and this is cleans local built image, on next lms startup dockerhub image will be downloaded
clean-lms:
	docker rmi krelms/janus-gateway:latest

build:
	@echo "$(OK_COLOR)==> Building default binary... $(NO_COLOR)"
	@CGO_ENABLED=0 go build -mod=vendor ${GO_LINKER_FLAGS} -o "dist/janus"

test-unit:
	@echo "$(OK_COLOR)==> Running unit tests$(NO_COLOR)"
	@go test ./...

test-integration: _mocks
	@echo "$(OK_COLOR)==> Running integration tests$(NO_COLOR)"
	@go test -cover -tags=integration -coverprofile=coverage.txt -covermode=atomic ./...

test-features: build _mocks
	@/bin/sh -c "./build/features.sh"

lint:
	@echo "$(OK_COLOR)==> Linting with golangci-lint running in docker container$(NO_COLOR)"
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.30.0 golangci-lint run -v

_mocks:
	@/bin/sh -c "./build/mocks.sh"
