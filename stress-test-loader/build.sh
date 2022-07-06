#!/usr/bin/env bash

#### build grpc_health_probe
protoc  --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative   proto/stresstestloader.proto

VERSION=`git describe --tags --long`

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -ldflags="-X main.Version=$(git describe --tags --long)" -o stress-test-loader-linux
gzip -f stress-test-loader-linux



