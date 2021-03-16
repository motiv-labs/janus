#!/bin/sh
cd /janus
if [ "$debug" == 1 ]; then
  echo "about to compile go for debugging"
  go build -gcflags "all=-N -l" -o main .
else
  echo "about to compile go"
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
fi
echo "compile finished"
if [ "$debug" == 1 ]; then
  dlv --listen=:40000 --headless=true --continue --accept-multiclient --api-version=2 exec ./main start
else
  ./main
fi
