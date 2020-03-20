NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# The import path is the unique absolute name of your repository.
# All subpackages should always be imported as relative to it.
# If you change this, run `make clean`.
PKG_SRC := github.com/hellofresh/janus

.PHONY: all clean deps build

all: clean deps test build

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/cucumber/godog/cmd/godog

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "JANUS_BUILD_ONLY_DEFAULT=$(JANUS_BUILD_ONLY_DEFAULT) PKG_SRC=$(PKG_SRC) VERSION=$(VERSION) ./build/build.sh"

test: lint format vet
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -v -cover ./...

test-integration: lint format vet
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -v -cover -tags=integration ./...

test-features:
	@/bin/sh -c "JANUS_BUILD_ONLY_DEFAULT=1 PKG_SRC=$(PKG_SRC) ./build/build.sh"
	@/bin/sh -c "./build/features.sh"

format:
	@echo "$(OK_COLOR)==> checking code formating with 'gofmt' tool$(NO_COLOR)"
	@gofmt -l -s cmd pkg | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

vet:
	@echo "$(OK_COLOR)==> checking code correctness with 'go vet' tool$(NO_COLOR)"
	@go vet ./...

lint: tools.golint
	@echo "$(OK_COLOR)==> checking code style with 'golint' tool$(NO_COLOR)"
	@go list ./... | xargs -n 1 golint -set_exit_status

clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@go clean
	@rm -rf bin $GOPATH/bin

#---------------
#-- tools
#---------------

.PHONY: tools tools.golint
tools: tools.golint

tools.golint:
	@command -v golint >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing golint"; \
		go get github.com/golang/lint/golint; \
	fi

