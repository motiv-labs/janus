#!/bin/bash

# Ensure this script fails if anything errors
set -e

# Extract the release artifact
tar -xzf release-candidate/*.tar.gz -C artifacts/

# Copy files needed for docker build
cp source-code/Dockerfile source-code/Makefile artifacts/

# Copy folders needed for docker build
cp -r source-code/docker artifacts/
