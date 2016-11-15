#!/bin/sh

# Ensure this script fails if anything errors
set -e

# Creates the necessary directories
mkdir -p docker-images/dev
mkdir -p docker-images/latest

# Copies code to the image directories
cp -r source-code docker-images/dev
cp -r source-code docker-images/latest

# Copies code to the image directories
cp docker-images/dev/ci/assets/Dockerfile.dev docker-images/dev/Dockerfile
echo "dev" >> docker-images/dev/version
cp docker-images/latest/ci/assets/Dockerfile docker-images/dev/Dockerfile 
