#!/usr/bin/env sh
# Builds the School Bundle once the Install, Type and Test Gates pass.
set -eu

cd -- "$(dirname -- "$0")"

verify_toolchain() {
	if ! command -v npm >/dev/null 2>&1; then
		echo "❌ npm not Found"
		return 1
	fi
	echo "✅ Toolchain Found"
}

install_packages() {
	if [ -f package-lock.json ]; then
		npm ci --silent
	else
		npm install --silent
	fi
	echo "✅ Packages Installed"
}

verify_types() {
	if ! npm run --silent build; then
		echo "❌ Type Check Failed"
		return 1
	fi
	echo "✅ Type Check Passed"
}

verify_tests() {
	if ! npm test --silent; then
		echo "❌ Tests Failed"
		return 1
	fi
	echo "✅ Tests Passed"
}

verify_toolchain
install_packages
verify_types
verify_tests
echo "✅ Bundle Built at ./dist"
