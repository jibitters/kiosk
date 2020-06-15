#!/usr/bin/env sh

VERSION=v0.2.0

env GOOS=freebsd GOARCH=amd64 go build -o kisok-freebsd-$VERSION ./cmd/kisok
env GOOS=linux GOARCH=amd64 go build -o kisok-linux-$VERSION ./cmd/kisok
env GOOS=darwin GOARCH=amd64 go build -o kisok-macos-$VERSION ./cmd/kisok
env GOOS=windows GOARCH=amd64 go build -o kisok-windows-$VERSION.exe ./cmd/kisok
