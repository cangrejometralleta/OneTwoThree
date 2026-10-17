#!/usr/bin/env sh
# Builds the School Jar once the Toolchain, Lint and Test Gates pass.
set -eu

cd -- "$(dirname -- "$0")"
JAR="target/school-1.0.0.jar"
LOWEST_JDK=21

find_javac() {
	if [ -n "${JAVA_HOME:-}" ]; then
		echo "$JAVA_HOME/bin/javac"
	else
		command -v javac 2>/dev/null || true
	fi
}

verify_toolchain() {
	javac_path="$(find_javac)"
	if [ -z "$javac_path" ] || [ ! -x "$javac_path" ]; then
		echo "❌ No JDK Found; a JRE Reports 'No compiler is provided'"
		echo "Set JAVA_HOME to a JDK $LOWEST_JDK or newer: JAVA_HOME=/usr/lib/jvm/<jdk> ./build.sh"
		return 1
	fi

	version="$("$javac_path" -version 2>&1 | sed 's/^javac \([0-9]*\).*/\1/')"
	if [ "$version" -lt "$LOWEST_JDK" ]; then
		echo "❌ javac $version is older than $LOWEST_JDK"
		echo "Set JAVA_HOME to a JDK $LOWEST_JDK or newer."
		return 1
	fi
	echo "✅ Toolchain Found: javac $version, Compiling for $LOWEST_JDK"
}

verify_tests() {
	# The JDK Ships no Formatter, so -Xlint:all -Werror in pom.xml is the Format Gate.
	# clean, because Maven's incremental Compiler Emits stale Records.
	if ! ./mvnw --quiet clean test; then
		echo "❌ Lint or Tests Failed"
		return 1
	fi
	echo "✅ Lint and Tests Passed"
}

build_jar() {
	./mvnw --quiet package -DskipTests
	echo "✅ Jar Built at ./$JAR"
}

verify_toolchain
verify_tests
build_jar
