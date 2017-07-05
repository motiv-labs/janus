NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# The import path is the unique absolute name of your repository.
# All subpackages should always be imported as relative to it.
# If you change this, run `make clean`.
IMPORT_PATH := github.com/hellofresh/janus
PKG_SRC := $(IMPORT_PATH)/cmd/janus

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

.PHONY: all clean deps build

all: clean deps build

deps:
	@echo "$(OK_COLOR)==> Installing glide dependencies$(NO_COLOR)"
	@go get -u github.com/Masterminds/glide
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/DATA-DOG/godog/cmd/godog
	@glide install

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "RUN_INTEGRATION=0 PKG_SRC=$(PKG_SRC) VERSION=${VERSION} ./build/build.sh"

test:
	@/bin/sh -c "./build/test.sh $(allpackages)"

test-integration:
	@/bin/sh -c "RUN_INTEGRATION=1 ./build/test.sh $(allpackages)"

lint:
	@echo "$(OK_COLOR)==> Linting... $(NO_COLOR)"
	@golint $(allpackages)

clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@go clean
	@rm -rf bin $GOPATH/bin

# cd into the GOPATH to workaround ./... not following symlinks
_allpackages = $(shell ( go list ./... 2>&1 1>&3 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) 1>&2 ) 3>&1 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)))

# memoize allpackages, so that it's executed only once and only if used
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)
