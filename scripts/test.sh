#!/usr/bin/env bash

go generate ./cmd/kiosk/main.go
go test ./cmd/... ./internal/... ./test/... -cover -failfast
