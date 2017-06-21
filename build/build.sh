#!/bin/bash

set -e

# Get rid of existing binaries
rm -f dist/janus*

# Check if VERSION variable set and not empty, otherwise set to default value
if [ -z "$VERSION" ]; then
  VERSION="0.0.1-dev"
fi
echo "Building application version $VERSION"

echo "Building default binary"
CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/janus" $PKG_SRC

echo "Building healthchecker binary"
CGO_ENABLED=0 go build -o "dist/healthchecker-linux-amd64" $HEALTH_SRC

if [ ! -z "${JANUS_BUILD_ONLY_DEFAULT}" ]; then
    echo "Only default binary was requested to build"
    exit 0
fi

# Build 386 amd64 binaries
OS_PLATFORM_ARG=(linux darwin windows freebsd openbsd)
OS_ARCH_ARG=(386 amd64)
for OS in ${OS_PLATFORM_ARG[@]}; do
  for ARCH in ${OS_ARCH_ARG[@]}; do
    echo "Building binary for $OS/$ARCH..."
    GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/janus_$OS-$ARCH" $PKG_SRC
  done
done

# Build arm binaries
OS_PLATFORM_ARG=(linux)
OS_ARCH_ARG=(arm arm64)
for OS in ${OS_PLATFORM_ARG[@]}; do
  for ARCH in ${OS_ARCH_ARG[@]}; do
    echo "Building binary for $OS/$ARCH..."
    GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/janus_$OS-$ARCH" $PKG_SRC
  done
done
