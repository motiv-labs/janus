#!/bin/sh
mkdir -p ${PROJECT_SRC}

cp -r . ${PROJECT_SRC}
cd ${PROJECT_SRC}

pwd
make deps
make test