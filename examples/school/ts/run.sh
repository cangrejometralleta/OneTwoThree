#!/usr/bin/env sh
# Runs School from the Bundle. Pass the Server Adapter; stdlib is the Default.
set -eu

cd -- "$(dirname -- "$0")"
SERVER="${1:-stdlib}"
APP_ENV="${APP_ENV:-development}"
export APP_ENV SERVER

verify_adapter() {
	case "$SERVER" in
		stdlib|node|express|fastify) ;;
		*)
			echo "❌ Unknown Adapter: $SERVER"
			echo "Usage: ./run.sh [stdlib|node|express|fastify]"
			return 1
			;;
	esac
}

verify_secret() {
	if [ -z "${TOKEN_SECRET:-}" ]; then
		echo "❌ TOKEN_SECRET is Missing"
		echo "The Deployment Supplies it; no File Key Holds it."
		echo "Try: TOKEN_SECRET=local-secret ./run.sh $SERVER"
		return 1
	fi
	echo "✅ Secret Supplied"
}

verify_bundle() {
	if [ ! -f dist/main.js ]; then
		echo "❌ Bundle Missing; Run ./build.sh first"
		return 1
	fi
	echo "✅ Bundle Found"
}

start_service() {
	echo "✅ Starting School on $SERVER in $APP_ENV"
	exec npm start --silent
}

verify_adapter
verify_secret
verify_bundle
start_service
