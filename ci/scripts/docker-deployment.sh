#!/bin/sh

# Ensure this script fails if anything errors
set -e

# Creates the necessary directories
mkdir -p docker-images
mkdir -p docker-images/dist

# unzip the binary
tar -zxf -C docker-images/dist release-candidate/janus.tar.gz

# Copies code to the image directories
cp docker-images/ci/assets/Dockerfile docker-images/Dockerfile
