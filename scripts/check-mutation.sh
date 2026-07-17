#!/usr/bin/env bash
set -euo pipefail

trap 'rm -f report.json' EXIT

go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@v0.0.0-20251226130216-48d0401f00fb \
	--exec-timeout 20 \
	value.go schedule.go exception.go query.go ranges.go timezone.go \
	search.go composition.go
