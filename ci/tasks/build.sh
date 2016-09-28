#!/bin/sh
set -x

CWD=$(pwd)
DEST=api-gateway

mkdir -p ${PROJECT_SRC}
cp -r . ${PROJECT_SRC}
cd ${PROJECT_SRC}
#build go binary
make

cd /go/bin
echo "* Creating tar.gz"
tar -czf ${DEST}.tar.gz ${DEST} > /dev/null
echo "* copying ${DEST}.tar.gz ${CWD}/artifacts"
cp ${DEST}.tar.gz ${CWD}/artifacts
