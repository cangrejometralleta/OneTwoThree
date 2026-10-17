package com.example;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Stream;

import org.junit.jupiter.api.Test;

/**
 * One Module cannot Let the Compiler Hold every Boundary, so this Test Reads the Imports.
 * It Says which Vendor may Live where, and the Build File is Part of the Check.
 */
class BoundariesTest {

  private static final Path SOURCES = Path.of("src", "main", "java", "com", "example");
  private static final Pattern IMPORT = Pattern.compile("^import (?:static )?([a-zA-Z0-9_.]+)", Pattern.MULTILINE);

  /** Each Vendor, and the one Package Allowed to Name it. */
  private static final Map<String, String> VENDOR_HOMES = Map.of(
      "com.fasterxml.jackson", "codec",
      "org.springframework", "serving/spring",
      "jakarta.", "serving/spring",
      "com.sun.net.httpserver", "serving/stdlib",
      "java.sql", "store",
      "org.h2", "nowhere: JDBC Finds the Driver at Runtime",
      "javax.crypto", "tokens");

  private static List<String[]> readImports() throws IOException {
    List<String[]> imports = new ArrayList<>();
    try (Stream<Path> files = Files.walk(SOURCES)) {
      for (Path file : files.filter(path -> path.toString().endsWith(".java")).toList()) {
        String home = SOURCES.relativize(file.getParent()).toString().replace('\\', '/');
        Matcher found = IMPORT.matcher(Files.readString(file));
        while (found.find()) {
          imports.add(new String[] {home, found.group(1), file.getFileName().toString()});
        }
      }
    }

    return imports;
  }

  @Test
  void everyVendorLivesInOnePackage() throws IOException {
    for (String[] line : readImports()) {
      VENDOR_HOMES.forEach((vendor, home) -> {
        if (line[1].startsWith(vendor)) {
          assertEquals(home, line[0], line[2] + " Imports " + line[1]);
        }
      });
    }
  }

  @Test
  void theCoreNamesNoVendorAndNoSocket() throws IOException {
    for (String[] line : readImports()) {
      if (line[0].equals("school")) {
        assertTrue(line[1].startsWith("java.util") || line[1].startsWith("com.example.school")
            || line[1].equals("com.example.faults.Fault"), line[2] + " Imports " + line[1]);
      }
    }
  }

  @Test
  void theShapesTheLayersSpeakInImportNothingOfOurs() throws IOException {
    for (String[] line : readImports()) {
      if (line[0].equals("transport") || line[0].equals("wire") || line[0].equals("faults")) {
        assertTrue(line[1].startsWith("java."), line[2] + " Imports " + line[1]);
      }
    }
  }

  @Test
  void theDriverIsOnlyEverFoundAtRuntime() throws IOException {
    String build = Files.readString(Path.of("pom.xml"));
    Matcher h2 = Pattern.compile("<artifactId>h2</artifactId>\\s*<scope>([a-z]+)</scope>").matcher(build);

    assertTrue(h2.find(), "pom.xml Must Declare the H2 Driver with a Scope");
    assertEquals("runtime", h2.group(1));
  }
}
