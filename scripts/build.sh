#!/usr/bin/env bash

KIOSK_VERSION=v0.0.9

go generate ./cmd/kiosk/main.go

env GOOS=linux GOARCH=amd64 go build -o kiosk-linux-$KIOSK_VERSION ./cmd/kiosk
env GOOS=freebsd GOARCH=amd64 go build -o kiosk-freebsd-$KIOSK_VERSION ./cmd/kiosk
env GOOS=darwin GOARCH=amd64 go build -o kiosk-macos-$KIOSK_VERSION ./cmd/kiosk
env GOOS=windows GOARCH=amd64 go build -o kiosk-windows-$KIOSK_VERSION.exe ./cmd/kiosk
