#!/usr/bin/env sh
# Builds the Converter once the Format, Vet and Test Gates pass.
set -eu

cd -- "$(dirname -- "$0")"
BINARY="bin/pdf"

verify_format() {
	unformatted=$(gofmt -l .)
	if [ -n "$unformatted" ]; then
		printf '%s\n' "$unformatted"
		echo "❌ Format Check Failed"
		return 1
	fi
	echo "✅ Format Check Passed"
}

verify_vet() {
	if ! go vet ./...; then
		echo "❌ Vet Failed"
		return 1
	fi
	echo "✅ Vet Passed"
}

verify_tests() {
	if ! go test ./...; then
		echo "❌ Tests Failed"
		return 1
	fi
	echo "✅ Tests Passed"
}

build_binary() {
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$BINARY" .
	echo "✅ Binary Built at ./$BINARY"
}

verify_format
verify_vet
verify_tests
build_binary
