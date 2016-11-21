NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

REPO=github.com/hellofresh/janus
GO_LINKER_FLAGS=-ldflags="-s -w"
GO_PROJECT_FILES=`go list -f '{{.Dir}}' ./... | grep -v /vendor/ | sed -n '1!p'`
GO_PROJECT_PACKAGES=`go list ./... | grep -v /vendor/`

# This how we want to name the binary output
DIR_OUT=$(CURDIR)/out
PROJECT_SRC=$(REPO)/cmd/janus
BINARY=janus

.PHONY: all clean deps build

all: clean deps build

deps:
	@echo "$(OK_COLOR)==> Installing glide dependencies$(NO_COLOR)"
	@curl https://glide.sh/get | sh
	@glide install

# Builds the project
build: build-ensure-dir build-linux build-osx

build-ensure-dir:
	@mkdir -p ${DIR_OUT}

build-linux:
	@echo "$(OK_COLOR)==> Building Linux amd64"
	@env GOOS=linux GOARCH=amd64 go build -o ${DIR_OUT}/${BINARY}-linux ${GO_LINKER_FLAGS} ${PROJECT_SRC}

build-osx:
	@echo "$(OK_COLOR)==> Building OSX amd64"
	@env GOOS=darwin GOARCH=amd64 go build -o ${DIR_OUT}/${BINARY}-darwin ${GO_LINKER_FLAGS} ${PROJECT_SRC}

# Installs our project: copies binaries
install:
	@echo "$(OK_COLOR)==> Installing project$(NO_COLOR)"
	go install -v

test:
	go test ${GO_PROJECT_PACKAGES} -v
	
# Cleans our project: deletes binaries
clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	if [ -d ${DIR_OUT} ] ; then rm -f ${DIR_OUT}/* ; fi
