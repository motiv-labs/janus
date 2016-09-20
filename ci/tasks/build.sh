#!/bin/sh
CWD=$(pwd)
DEST=api-gateway
cd ${PROJECT_SRC}

make

cd /go/bin

echo "* Creating tar.gz"
tar -cvzf ${DEST}.tar.gz ${DEST} > /dev/null

cp ${DEST}.tar.gz ${CWD}/artifacts
