#!/usr/bin/env bash

go generate ./cmd/kiosk/main.go

set -e
echo "" > coverage.txt

for d in $(go list ./...); do
    go test -coverprofile=profile.out -failfast "$d"
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
