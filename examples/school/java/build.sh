#!/usr/bin/env sh
# Builds the School Jar once the Toolchain and Test Gates pass.
set -eu

cd -- "$(dirname -- "$0")"
JAR="target/school-1.0.0.jar"

verify_toolchain() {
	if [ -n "${JAVA_HOME:-}" ] && [ ! -x "$JAVA_HOME/bin/javac" ]; then
		echo "❌ JAVA_HOME Points at a JRE: $JAVA_HOME"
		echo "The Build Wants a JDK; a JRE Reports 'No compiler is provided'."
		return 1
	fi
	if [ -z "${JAVA_HOME:-}" ] && ! command -v javac >/dev/null 2>&1; then
		echo "❌ No javac Found"
		echo "Set JAVA_HOME to a JDK, or Put javac on the Path."
		return 1
	fi
	echo "✅ Toolchain Found"
}

verify_tests() {
	# clean, because Maven's incremental Compiler Emits stale Records
	# and an Editor's Language Server Writes into the same target.
	if ! ./mvnw --quiet clean test; then
		echo "❌ Tests Failed"
		return 1
	fi
	echo "✅ Tests Passed"
}

build_jar() {
	./mvnw --quiet package -DskipTests
	echo "✅ Jar Built at ./$JAR"
}

verify_toolchain
verify_tests
build_jar
