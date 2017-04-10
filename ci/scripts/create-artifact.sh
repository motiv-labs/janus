#!/bin/sh

# Ensure this script fails if anything errors
set -e

CWD=$(pwd)
BINARY=janus

VERSION=$(cat ./version/version | sed 's/-.*$//g')
cd ./source-code/ && GITHASH=$(git rev-parse --short HEAD) && cd ${CWD}

# Goes to the application source code
cd source-code/
# Creates the project source on the gopath
mkdir -p ${PROJECT_SRC}
# Copies the current source code from the app to the gopath
cp -r . ${PROJECT_SRC}
# Goes to the application on the gopath
cd ${PROJECT_SRC}
# Set version to inject it into application in compile-time
export VERSION="${VERSION}-${GITHASH}"
# Build the go application
make

# Goes to the generated go binaries
cd dist

# Pack binaries
for i in ./*; do
    echo "Packing binary for $i..."
    tar -czf $i.tar.gz $i
done

# Copies the tar to the artifact folder so its available to the next step of the pipeline
echo "Copying *.tar.gz ${CWD}/artifacts"
cp *.tar.gz ${CWD}/artifacts
