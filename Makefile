NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

.PHONY: all lint test-unit test-integration test-features build

all: test-unit build

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "JANUS_BUILD_ONLY_DEFAULT=$(JANUS_BUILD_ONLY_DEFAULT) VERSION=$(VERSION) ./build/build.sh"

test-unit:
	@echo "$(OK_COLOR)==> Running unit tests$(NO_COLOR)"
	@go test ./...

test-integration: _mocks
	@echo "$(OK_COLOR)==> Running integration tests$(NO_COLOR)"
	@go test -cover -tags=integration -coverprofile=coverage.txt -covermode=atomic ./...

test-features: _mocks
	@/bin/sh -c "JANUS_BUILD_ONLY_DEFAULT=1 ./build/build.sh"
	@/bin/sh -c "./build/features.sh"

lint:
	@echo "$(OK_COLOR)==> Linting with golangci-lint running in docker container$(NO_COLOR)"
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.30.0 golangci-lint run -v

_mocks:
	@/bin/sh -c "./build/mocks.sh"
