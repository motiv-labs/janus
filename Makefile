NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# The import path is the unique absolute name of your repository.
# All subpackages should always be imported as relative to it.
# If you change this, run `make clean`.
PKG_SRC := github.com/hellofresh/janus

.PHONY: all clean deps lint test build

all: clean deps lint test build

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -u github.com/cucumber/godog/cmd/godog@v0.10.0

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "JANUS_BUILD_ONLY_DEFAULT=$(JANUS_BUILD_ONLY_DEFAULT) PKG_SRC=$(PKG_SRC) VERSION=$(VERSION) ./build/build.sh"

test: lint
	@echo "$(OK_COLOR)==> Running unit tests$(NO_COLOR)"
	@go test -cover ./...

test-integration: _mocks
	@echo "$(OK_COLOR)==> Running integration tests$(NO_COLOR)"
	@go test -cover -tags=integration ./...

test-features: _mocks
	@/bin/sh -c "JANUS_BUILD_ONLY_DEFAULT=1 PKG_SRC=$(PKG_SRC) ./build/build.sh"
	@/bin/sh -c "./build/features.sh"

lint:
	@echo "$(OK_COLOR)==> Linting with golangci-lint$(NO_COLOR)"
	@golangci-lint run

lint-docker:
	@echo "$(OK_COLOR)==> Linting with golangci-lint running in docker container$(NO_COLOR)"
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.30.0 golangci-lint run -v

clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@go clean
	@rm -rf bin $GOPATH/bin

_mocks:
	@/bin/sh -c "./build/mocks.sh"
