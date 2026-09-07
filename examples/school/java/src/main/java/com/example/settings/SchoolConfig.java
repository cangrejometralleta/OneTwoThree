package com.example.settings;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.node.IntNode;
import com.fasterxml.jackson.databind.node.TextNode;
import com.example.settings.Settings.ValueRule;

import java.nio.file.Path;
import java.util.LinkedHashMap;
import java.util.Map;

/** SchoolConfig Holds Deployment Values; Secrets Have no File Key. */
@JsonIgnoreProperties(ignoreUnknown = false)
public record SchoolConfig(int port, String databasePath, int tokenLifeSeconds, String serverAdapter,
    String tokenSecret) {

  /** The Bounds a Deployment Value Must Fall inside. */
  private static final int LOWEST_PORT = 1;
  private static final int HIGHEST_PORT = 65535;
  private static final int SECONDS_IN_A_DAY = 60 * 60 * 24;

  private static final Map<String, ValueRule> RULES = Map.of(
      "port", Settings.checkIntegerRange(LOWEST_PORT, HIGHEST_PORT),
      "databasePath", Settings.CHECK_TEXT_VALUE,
      "tokenLifeSeconds", Settings.checkIntegerRange(1, SECONDS_IN_A_DAY),
      "serverAdapter", Settings.checkChoiceValue("stdlib"));

  private static final Map<String, String> NAMES = Map.of(
      "SCHOOL_PORT", "port",
      "SCHOOL_DATABASE_PATH", "databasePath",
      "SCHOOL_TOKEN_LIFE_SECONDS", "tokenLifeSeconds",
      "SERVER", "serverAdapter");

  /** LoadSchoolConfig Resolves Startup Inputs before any Store or Listener Opens. */
  public static SchoolConfig loadSchoolConfig(Path root, Map<String, String> environment) {
    checkEnvironmentKeys(environment);

    Map<String, JsonNode> values = readEnvironmentFiles(root, environment.get("APP_ENV"));
    applyEnvironmentValues(values, environment);

    Deployment deployment = Settings.decodeDataValues(values, RULES, Deployment.class);

    String secret = environment.getOrDefault("TOKEN_SECRET", "");
    if (secret.isBlank()) {
      throw new IllegalStateException("TOKEN_SECRET is Missing");
    }

    return new SchoolConfig(deployment.port(), deployment.databasePath(), deployment.tokenLifeSeconds(),
        deployment.serverAdapter(), secret);
  }

  private record Deployment(int port, String databasePath, int tokenLifeSeconds, String serverAdapter) {
  }

  /** ReadEnvironmentFiles Applies Defaults before the selected Environment. */
  private static Map<String, JsonNode> readEnvironmentFiles(Path root, String environment) {
    if (!"development".equals(environment) && !"production".equals(environment)) {
      throw new IllegalStateException("APP_ENV Must be development or production");
    }

    Path folder = root.resolve("config");
    Map<String, JsonNode> values = new LinkedHashMap<>(
        Settings.readDataFile(folder.resolve("school.defaults.json"), RULES));
    values.putAll(Settings.readDataFile(folder.resolve("school." + environment + ".json"), RULES));

    return values;
  }

  /** CheckEnvironmentKeys Refuses undeclared Variables in this Service Namespace. */
  private static void checkEnvironmentKeys(Map<String, String> environment) {
    for (String name : environment.keySet()) {
      if (name.startsWith("SCHOOL_") && !NAMES.containsKey(name)) {
        throw new IllegalStateException("Unknown Variable " + name + "; Global Constants Cannot be Overridden");
      }
    }
  }

  /** ApplyEnvironmentValues Validates declared Variables before Replacing Values. */
  private static void applyEnvironmentValues(Map<String, JsonNode> values, Map<String, String> environment) {
    NAMES.forEach((name, key) -> {
      String value = environment.get(name);
      if (value == null) {
        return;
      }
      try {
        values.put(key, encodeEnvironmentValue(value, RULES.get(key)));
      } catch (RuntimeException broken) {
        throw new IllegalStateException(name + ": " + broken.getMessage(), broken);
      }
    });
  }

  /** EncodeEnvironmentValue Accepts Text or a decimal Integer, never implicit Booleans. */
  private static JsonNode encodeEnvironmentValue(String value, ValueRule rule) {
    JsonNode text = TextNode.valueOf(value);
    try {
      rule.accept(text);
      return text;
    } catch (IllegalArgumentException notText) {
      if (!value.matches("^[0-9]+$")) {
        throw new IllegalArgumentException("Invalid Value");
      }
    }

    JsonNode number = IntNode.valueOf(Integer.parseInt(value));
    rule.accept(number);

    return number;
  }
}
