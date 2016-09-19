NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# This how we want to name the binary output
BINARY=auth-service

.PHONY: all clean deps install

all: clean deps install

deps: 
	@echo "$(OK_COLOR)==> Installing gin proxy$(NO_COLOR)"
	@go get github.com/codegangsta/gin

	@echo "$(OK_COLOR)==> Installing glide dependencies$(NO_COLOR)"
	@curl https://glide.sh/get | sh
	@glide install

# Builds the project
build:
	@echo "$(OK_COLOR)==> Building project$(NO_COLOR)"
	@go build -o ${BINARY}

# Installs our project: copies binaries
install:
	@echo "$(OK_COLOR)==> Installing project$(NO_COLOR)"
	go install -v

test:
	go test -v
	
# Cleans our project: deletes binaries
clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
