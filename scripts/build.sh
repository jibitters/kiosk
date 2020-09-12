#!/usr/bin/env sh

VERSION=v1.0.4

env GOOS=freebsd GOARCH=amd64 go build -o kiosk-freebsd-$VERSION ./cmd/kiosk
env GOOS=linux GOARCH=amd64 go build -o kiosk-linux-$VERSION ./cmd/kiosk
env GOOS=darwin GOARCH=amd64 go build -o kiosk-macos-$VERSION ./cmd/kiosk
env GOOS=windows GOARCH=amd64 go build -o kiosk-windows-$VERSION.exe ./cmd/kiosk
