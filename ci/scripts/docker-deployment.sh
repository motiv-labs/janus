#!/bin/sh

# Ensure this script fails if anything errors
set -e

# Creates the necessary directories
mkdir -p docker-images
mkdir -p docker-images/dist
mkdir -p docker-images/ci/assets

# unzip the binary
tar -C docker-images/dist -zxf release-candidate/janus.tar.gz

# Copies code to the image directories
cp source-code/ci/assets/* docker-images/ci/assets
cp source-code/Dockerfile docker-images/Dockerfile
