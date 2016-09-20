#!/bin/sh
CWD=$(pwd)
cd ${PROJECT_SRC}

make
cp /go/bin/api-gateway ${CWD}/artifacts