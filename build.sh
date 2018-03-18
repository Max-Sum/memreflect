#!/bin/sh

go get -t github.com/max-sum/memreflect/build
GOOS=linux GOARCH=amd64 go build -o memreflect github.com/max-sum/memreflect/build
docker build .