package com.example.settings;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

class SchoolConfigTest {

  private static final Set<String> ADAPTERS = Set.of("stdlib", "spring");
  private static final Path SHIPPED = Path.of("..");

  @TempDir
  Path root;

  /** buildConfigFixture Keeps broken Files away from the checked-in Examples. */
  @BeforeEach
  void buildConfigFixture() throws IOException {
    Files.createDirectories(root.resolve("config"));
    Files.createDirectories(root.resolve("constants"));
    for (String name : List.of("school.defaults.json", "school.development.json", "school.production.json",
        "registry.json")) {
      Files.copy(SHIPPED.resolve("config").resolve(name), root.resolve("config").resolve(name));
    }
    Files.copy(SHIPPED.resolve("constants").resolve("school.json"), root.resolve("constants").resolve("school.json"));
  }

  private static Map<String, String> readyEnvironment() {
    Map<String, String> environment = new HashMap<>();
    environment.put("APP_ENV", "production");
    environment.put("TOKEN_SECRET", "test-secret");

    return environment;
  }

  private SchoolConfig load(Map<String, String> environment) throws SettingsRefused {
    return SchoolConfig.loadSchoolConfig(root, environment, ADAPTERS);
  }

  private void writeProduction(String content) throws IOException {
    Files.writeString(root.resolve("config").resolve("school.production.json"), content);
  }

  @Test
  void theShippedFilesLoadForBothEnvironments() throws SettingsRefused {
    Map<String, String> environment = readyEnvironment();
    assertEquals(Path.of("school.production.db"), load(environment).databasePath());

    environment.put("APP_ENV", "development");
    SchoolConfig development = load(environment);
    assertEquals(Path.of("school.db"), development.databasePath());
    assertEquals(8080, development.port());
    assertEquals(600, development.tokenLifeSeconds());
    assertEquals("stdlib", development.serverAdapter());
  }

  @Test
  void defaultsThenTheEnvironmentFileThenTheVariables() throws Exception {
    writeProduction("{\"databasePath\":\"production.db\",\"port\":9000}");
    Map<String, String> environment = readyEnvironment();

    SchoolConfig fromFile = load(environment);
    assertEquals(9000, fromFile.port());
    assertEquals(Path.of("production.db"), fromFile.databasePath());

    environment.put("SCHOOL_PORT", "9100");
    environment.put("SERVER", "spring");
    SchoolConfig fromVariables = load(environment);
    assertEquals(9100, fromVariables.port());
    assertEquals("spring", fromVariables.serverAdapter());
  }

  @Test
  void theRegistryResolvesFromTheDataRootAndMustExist() throws SettingsRefused {
    Map<String, String> environment = readyEnvironment();
    assertEquals(root.resolve("config").resolve("registry.json"), load(environment).registryPath());

    environment.put("SCHOOL_REGISTRY_PATH", root.resolve("config").resolve("absent.json").toString());
    assertThrows(SettingsRefused.class, () -> load(environment));

    environment.put("SCHOOL_REGISTRY_PATH", "config");
    assertThrows(SettingsRefused.class, () -> load(environment));
  }

  @Test
  void startupRefusesWhatItCannotTrustInTheEnvironment() {
    Map<String, Map<String, String>> cases = new HashMap<>();
    cases.put("missing environment", Map.of("APP_ENV", ""));
    cases.put("unknown environment", Map.of("APP_ENV", "staging"));
    cases.put("path traversal", Map.of("APP_ENV", "../production"));
    cases.put("empty override", Map.of("SCHOOL_PORT", ""));
    cases.put("signed port", Map.of("SCHOOL_PORT", "+80"));
    cases.put("spaced port", Map.of("SCHOOL_PORT", " 80"));
    cases.put("fractional port", Map.of("SCHOOL_PORT", "80.0"));
    cases.put("port too high", Map.of("SCHOOL_PORT", "65536"));
    cases.put("port zero", Map.of("SCHOOL_PORT", "0"));
    cases.put("empty path", Map.of("SCHOOL_DATABASE_PATH", " "));
    cases.put("empty registry path", Map.of("SCHOOL_REGISTRY_PATH", ""));
    cases.put("global variable", Map.of("SCHOOL_MINIMUM_AGE_YEARS", "1"));
    cases.put("unknown variable", Map.of("SCHOOL_TYPO", "1"));
    cases.put("missing secret", Map.of("TOKEN_SECRET", ""));
    cases.put("unknown adapter", Map.of("SERVER", "gin"));
    cases.put("empty adapter", Map.of("SERVER", ""));
    cases.put("lifetime zero", Map.of("SCHOOL_TOKEN_LIFE_SECONDS", "0"));
    cases.put("lifetime beyond a day", Map.of("SCHOOL_TOKEN_LIFE_SECONDS", "86401"));

    cases.forEach((name, overrides) -> {
      Map<String, String> environment = readyEnvironment();
      environment.putAll(overrides);
      assertThrows(SettingsRefused.class, () -> load(environment), name);
    });

    Map<String, String> unsecret = readyEnvironment();
    unsecret.remove("TOKEN_SECRET");
    assertThrows(SettingsRefused.class, () -> load(unsecret));
  }

  @Test
  void startupRefusesWhatItCannotTrustInTheFiles() throws IOException {
    Map<String, String> cases = new HashMap<>();
    cases.put("global file key", "{\"databasePath\":\"test.db\",\"minimumAgeYears\":1}");
    cases.put("unknown file key", "{\"databasePath\":\"test.db\",\"typo\":1}");
    cases.put("null key", "{\"databasePath\":null}");
    cases.put("invalid type", "{\"databasePath\":\"test.db\",\"port\":\"9000\"}");
    cases.put("fractional", "{\"databasePath\":\"test.db\",\"port\":9000.5}");
    cases.put("invalid file range", "{\"databasePath\":\"test.db\",\"port\":0}");
    cases.put("invalid json", "{");
    cases.put("null document", "null");
    cases.put("array document", "[]");
    cases.put("empty file", "");
    cases.put("secret in file", "{\"databasePath\":\"test.db\",\"tokenSecret\":\"forbidden\"}");

    for (Map.Entry<String, String> file : cases.entrySet()) {
      writeProduction(file.getValue());
      assertThrows(SettingsRefused.class, () -> load(readyEnvironment()), file.getKey());
    }

    Files.delete(root.resolve("config").resolve("school.production.json"));
    assertThrows(SettingsRefused.class, () -> load(readyEnvironment()));
  }

  @Test
  void aMissingKeyStopsStartup() throws IOException {
    Files.writeString(root.resolve("config").resolve("school.defaults.json"), "{\"port\":8080}");
    writeProduction("{}");

    assertThrows(SettingsRefused.class, () -> load(readyEnvironment()));
  }

  @Test
  void anOverrideNeverHidesAnInvalidDefault() throws IOException {
    Files.writeString(root.resolve("config").resolve("school.defaults.json"), "{\"port\":null}");
    Map<String, String> environment = readyEnvironment();
    environment.put("SCHOOL_PORT", "9100");

    assertThrows(SettingsRefused.class, () -> load(environment));
  }

  @Test
  void aVariableOutsideTheNamespaceIsIgnored() throws SettingsRefused {
    Map<String, String> environment = readyEnvironment();
    environment.put("SCHOOLING", "anything");
    environment.put("PATH", "/usr/bin");

    assertEquals(8080, load(environment).port());
  }

  @Test
  void theSecretNeverPrints() throws SettingsRefused {
    assertFalse(load(readyEnvironment()).toString().contains("test-secret"));
  }

  @Test
  void theGlobalConstantsHoldTheEnrolmentAgeAndRefuseBrokenData() throws Exception {
    assertEquals(18, Settings.readGlobalConstants(root).minimumAgeYears());

    for (String broken : List.of("{\"minimumAgeYears\":null}", "{}", "{\"minimumAgeYears\":0}",
        "{\"minimumAgeYears\":18,\"port\":1}")) {
      Files.writeString(root.resolve("constants").resolve("school.json"), broken);
      assertThrows(SettingsRefused.class, () -> Settings.readGlobalConstants(root), broken);
    }
  }

  @Test
  void theDataRootIsFoundAboveTheWorkingDirectory() throws SettingsRefused {
    assertEquals(SHIPPED.toAbsolutePath().normalize(),
        Settings.findSchoolDataRoot(Path.of("")).toAbsolutePath().normalize());
  }
}
