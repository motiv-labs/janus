#!/bin/sh
CWD=$(pwd)

# Defines the zip name
DEST=janus

# Goes to the application source code
cd app-master

# Creates the project source on the gopath
mkdir -p ${PROJECT_SRC}

# Copies the current srouce code from the app to the gopath
cp -r . ${PROJECT_SRC}

# Goes to the application on the gopath
cd ${PROJECT_SRC}

# Build the go application
make

# Goes to the generated go binaries
cd /go/bin

# Zip all binaries in one sigle tar
echo "* Creating tar.gz"
tar -czvf ${DEST}.tar.gz ${DEST}

# Copies the tar to the artifact folder so its available to the next step of the pipeline
echo "* Copying *.tar.gz ${CWD}/artifacts"
cp *.tar.gz ${CWD}/artifacts
