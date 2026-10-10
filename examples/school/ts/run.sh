#!/usr/bin/env sh
set -eu
cd -- "$(dirname -- "$0")"

verify_node() {
	if ! command -v node >/dev/null 2>&1; then
		echo "Node.js 22.6 or newer is required" >&2
		exit 1
	fi
	node -e 'const [major,minor]=process.versions.node.split(".").map(Number); if(major<22 || (major===22 && minor<6)){console.error("Node.js 22.6 or newer is required");process.exit(1)}'
}

start_service() {
	export DATA_ROOT="${DATA_ROOT:-..}"
	export SCHOOL_DATABASE_PATH="${SCHOOL_DATABASE_PATH:-school-ts.db}"
	exec node --experimental-strip-types main.ts
}

verify_node
start_service
