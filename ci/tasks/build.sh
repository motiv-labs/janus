#!/bin/sh

CWD=$(pwd)
DEST=api-gateway

mkdir -p ${PROJECT_SRC}
cp -r . ${PROJECT_SRC}
cd ${PROJECT_SRC}

make

cd /go/bin

echo "* Creating tar.gz"
tar -czf ${DEST}.tar.gz ${DEST} > /dev/null

cp ${DEST}.tar.gz ${CWD}/artifacts
