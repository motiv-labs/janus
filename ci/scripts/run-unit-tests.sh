#!/bin/sh

# Ensure this script fails if anything errors
set -e

# Creates the project source on the gopath
mkdir -p ${PROJECT_SRC}

# Copies the current source code from the app to the gopath
cp -r . ${PROJECT_SRC}

# Goes to the application on the gopath
cd ${PROJECT_SRC}

make deps
make test
