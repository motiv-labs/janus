#!/bin/sh

# Ensure this script fails if anything errors
set -e

# Creates the necessary directories
mkdir -p docker-images/latest

# Copies code to the image directories
cp -r "source-code/"* docker-images/latest

# Copies code to the image directories
cp docker-images/latest/ci/assets/Dockerfile docker-images/latest/Dockerfile
