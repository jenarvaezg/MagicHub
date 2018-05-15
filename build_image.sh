#!/bin/bash

export CGO_ENABLED=0
export GOOS=linux
export MAGICHUB_VERSION=$(cat version.txt)

go build -a -installsuffix cgo -o MagicHub .
docker build -t jenarvaezg/magichub:$MAGICHUB_VERSION .

