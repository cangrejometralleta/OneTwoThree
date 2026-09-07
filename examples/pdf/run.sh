#!/usr/bin/env sh
# Converts one Markdown Source into a PDF beside it.
set -eu

cd -- "$(dirname -- "$0")"
SOURCE="${1:-}"

verify_source() {
	if [ -z "$SOURCE" ]; then
		echo "❌ No Source Given"
		echo "Usage: ./run.sh <source.md>"
		echo "Try: ./run.sh testdata/sample.md"
		return 1
	fi

	if [ ! -f "$SOURCE" ]; then
		echo "❌ No such Source: $SOURCE"
		return 1
	fi
	echo "✅ Source Found"
}

convert_source() {
	exec go run . "$SOURCE"
}

verify_source
convert_source
