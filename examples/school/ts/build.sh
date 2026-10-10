#!/usr/bin/env sh
set -eu
cd -- "$(dirname -- "$0")"

verify_toolchain() {
	if ! command -v node >/dev/null 2>&1; then
		echo "Node.js 22.6 or newer is required" >&2
		exit 1
	fi
	node -e 'const [major,minor]=process.versions.node.split(".").map(Number); if(major<22 || (major===22 && minor<6)){console.error("Node.js 22.6 or newer is required");process.exit(1)}'
	if [ ! -x node_modules/.bin/tsc ]; then
		echo "TypeScript dependencies are missing; run npm install" >&2
		exit 1
	fi
}

verify_source() {
	node --experimental-strip-types --check main.ts
	node_modules/.bin/tsc --project tsconfig.json
}

run_tests() {
	node --experimental-strip-types tests/school.test.ts
}

verify_toolchain
verify_source
run_tests
