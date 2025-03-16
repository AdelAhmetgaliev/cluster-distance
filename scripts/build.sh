#!/usr/bin/env bash

gofmt -w .

GOOS=linux GOARCH=amd64 \
go build -buildmode=pie \
    -ldflags="-linkmode=external -s -w -bindnow" \
    -o ./bin/cluster-distance ./cmd/cluster-distance/main.go

