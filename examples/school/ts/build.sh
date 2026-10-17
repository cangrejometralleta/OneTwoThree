#!/usr/bin/env sh
# Checks School TypeScript once the Toolchain, Type and Test Gates pass.
set -eu

cd -- "$(dirname -- "$0")"
LOWEST_NODE=24

verify_toolchain() {
	if ! command -v node >/dev/null 2>&1; then
		echo "❌ No Node Found"
		echo "Install Node.js $LOWEST_NODE or newer: https://nodejs.org/en/download"
		return 1
	fi
	major="$(node -p 'process.versions.node.split(".")[0]')"
	if [ "$major" -lt "$LOWEST_NODE" ]; then
		echo "❌ Node $major is older than $LOWEST_NODE"
		echo "Node $LOWEST_NODE Strips Types and Ships the Test Runner this Build Uses."
		return 1
	fi
	echo "✅ Toolchain Found: Node $major"
}

install_packages() {
	if ! npm ci --no-audit --no-fund --silent; then
		echo "❌ Packages Failed to Install"
		echo "Run npm ci here; better-sqlite3 Builds a Native Binding."
		return 1
	fi
	echo "✅ Packages Installed"
}

verify_types() {
	# Node Ships no Formatter, so the strict Type Check is the Format Gate.
	if ! node_modules/.bin/tsc --project tsconfig.json; then
		echo "❌ Type Check Failed"
		return 1
	fi
	echo "✅ Type Check Passed"
}

verify_tests() {
	if ! node --test 'src/**/*.test.ts' >/dev/null; then
		echo "❌ Tests Failed; Run: node --test 'src/**/*.test.ts'"
		return 1
	fi
	echo "✅ Tests Passed"
}

verify_toolchain
install_packages
verify_types
verify_tests
