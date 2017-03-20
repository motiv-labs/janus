NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# The binary to build (just the basename).
BIN := hellowork-api

# This repo's root import path (under GOPATH).
PKG := github.com/hellofresh/janus

# Which architecture to build - see $(ALL_ARCH) for options.
ARCH ?= amd64

###
### These variables should not need tweaking.
###

SRC_DIRS := cmd pkg # directories which hold app source (not vendored)

.PHONY: all clean deps build

all: clean deps build

deps:
	@echo "$(OK_COLOR)==> Installing glide dependencies$(NO_COLOR)"
	@go get -u github.com/Masterminds/glide
	@glide install

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@/bin/sh -c "ARCH=$(ARCH) ./build/build.sh"

test:
	@/bin/sh -c "./build/test.sh $(SRC_DIRS)"

clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@go clean
