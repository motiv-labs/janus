#!/bin/sh

# Ensure this script fails if anything errors
set -e

CWD=$(pwd)
BINARY=janus

# Goes to the application source code
cd source-code/
# Creates the project source on the gopath
mkdir -p ${PROJECT_SRC}
# Copies the current source code from the app to the gopath
cp -r . ${PROJECT_SRC}
# Goes to the application on the gopath
cd ${PROJECT_SRC}
# Build the go application
make

# Goes to the generated go binaries
cd dist

# Pack 386 amd64 binaries
OS_PLATFORM_ARG=(linux darwin windows freebsd openbsd)
OS_ARCH_ARG=(386 amd64)
for OS in ${OS_PLATFORM_ARG[@]}; do
  for ARCH in ${OS_ARCH_ARG[@]}; do
    echo "Packing binary for $OS/$ARCH..."
    tar -czf $OS_$ARCH.tar.gz $BINARY_$OS-$ARCH
  done
done


# Pack arm binaries
OS_PLATFORM_ARG=(linux)
OS_ARCH_ARG=(arm arm64)
for OS in ${OS_PLATFORM_ARG[@]}; do
  for ARCH in ${OS_ARCH_ARG[@]}; do
    echo "Packing binary for $OS/$ARCH..."
    tar -czf $OS_$ARCH.tar.gz $BINARY_$OS-$ARCH
  done
done

# Copies the tar to the artifact folder so its available to the next step of the pipeline
echo "Copying *.tar.gz ${CWD}/artifacts"
cp *.tar.gz ${CWD}/artifacts
