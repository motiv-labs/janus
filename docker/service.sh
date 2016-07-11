#!/bin/bash

cd $GOPATH/src/github.com/hellofresh/api-gateway
exec /sbin/setuser api-gateway $GOPATH/bin/api-gateway
