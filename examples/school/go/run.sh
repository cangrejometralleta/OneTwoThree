#!/usr/bin/env sh
# Runs School. Pass the Server Adapter; stdlib is the Default.
set -eu

cd -- "$(dirname -- "$0")"
ADAPTER="${1:-}"

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
		echo "Copy ../.env.example to .env, or Export it."
		return 1
	fi
	echo "✅ Secret Supplied"
}



start_service() {
	echo "✅ Starting School on $SERVER in $APP_ENV"
	exec go run .
}

read_dotenv
select_settings
verify_adapter
verify_secret
start_service
