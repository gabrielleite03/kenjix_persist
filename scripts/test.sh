#!/usr/bin/env bash
set -euo pipefail

COVER_OUT=${1:-coverage.out}
HTML_OUT=${2:-coverage.html}

echo "Running go tests with coverage..."
go test ./... -coverprofile="$COVER_OUT"

echo "Coverage summary:"
go tool cover -func="$COVER_OUT"

echo "Generating HTML report: $HTML_OUT"
go tool cover -html="$COVER_OUT" -o "$HTML_OUT"

echo "Done. Coverage HTML: $HTML_OUT"
