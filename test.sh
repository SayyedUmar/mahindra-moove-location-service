#!/usr/bin/env bash
set -e
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
