#!/usr/bin/env sh
# Runs School from the Jar. Pass the Server Adapter; stdlib is the Default.
set -eu

cd -- "$(dirname -- "$0")"
ADAPTER="${1:-}"
JAR="target/school-1.0.0.jar"

read_dotenv() {
	if [ ! -f .env ]; then
		return 0
	fi

	while IFS='=' read -r name value; do
		case "$name" in ''|\#*) continue ;; esac
		if [ -z "$(printenv "$name" 2>/dev/null || true)" ]; then
			export "$name=$value"
		fi
	done < .env

	echo "✅ Loaded .env"
}

select_settings() {
	# The Argument Wins, then the Environment, then the Default.
	SERVER="${ADAPTER:-${SERVER:-stdlib}}"
	APP_ENV="${APP_ENV:-development}"
	export APP_ENV SERVER
}

verify_adapter() {
	case "$SERVER" in
		stdlib) ;;
		*)
			echo "❌ Unknown Adapter: $SERVER"
			echo "This Runtime Ships one Server. Usage: ./run.sh [stdlib]"
			return 1
			;;
	esac
}

verify_secret() {
	if [ -z "${TOKEN_SECRET:-}" ]; then
		echo "❌ TOKEN_SECRET is Missing"
		echo "The Deployment Supplies it; no File Key Holds it."
		echo "Copy ../.env.example to .env, or Export it."
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

read_dotenv
select_settings
verify_adapter
verify_secret
verify_jar
start_service
