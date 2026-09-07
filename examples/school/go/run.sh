#!/usr/bin/env sh
# Runs School from Source. Pass the Server Adapter; stdlib is the Default.
set -eu

cd -- "$(dirname -- "$0")"
SERVER="${1:-stdlib}"
APP_ENV="${APP_ENV:-development}"
export APP_ENV SERVER

verify_adapter() {
	case "$SERVER" in
		stdlib|chi|gin) ;;
		*)
			echo "❌ Unknown Adapter: $SERVER"
			echo "Usage: ./run.sh [stdlib|chi|gin]"
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

start_service() {
	echo "✅ Starting School on $SERVER in $APP_ENV"
	exec go run .
}

verify_adapter
verify_secret
start_service
