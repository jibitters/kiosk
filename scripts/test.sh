#!/usr/bin/env bash

go generate ./cmd/kiosk/main.go
go test ./... -coverpkg=all
