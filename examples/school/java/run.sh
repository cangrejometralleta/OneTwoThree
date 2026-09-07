#!/usr/bin/env sh
# Runs School from the Jar. Pass the Server Adapter; stdlib is the Default.
set -eu

cd -- "$(dirname -- "$0")"
SERVER="${1:-stdlib}"
APP_ENV="${APP_ENV:-development}"
JAR="target/school-1.0.0.jar"
export APP_ENV SERVER

verify_adapter() {
	if [ "$SERVER" != "stdlib" ]; then
		echo "❌ Unknown Adapter: $SERVER"
		echo "This Runtime Ships one Server. Usage: ./run.sh [stdlib]"
		return 1
	fi
}

verify_secret() {
	if [ -z "${TOKEN_SECRET:-}" ]; then
		echo "❌ TOKEN_SECRET is Missing"
		echo "The Deployment Supplies it; no File Key Holds it."
		echo "Try: TOKEN_SECRET=local-secret ./run.sh"
		return 1
	fi
	echo "✅ Secret Supplied"
}

verify_jar() {
	if [ ! -f "$JAR" ]; then
		echo "❌ Jar Missing; Run ./build.sh first"
		return 1
	fi
	echo "✅ Jar Found"
}

start_service() {
	echo "✅ Starting School on $SERVER in $APP_ENV"
	exec java -jar "$JAR"
}

verify_adapter
verify_secret
verify_jar
start_service
