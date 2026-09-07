package com.example.settings;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.function.Consumer;

/** The strict Loader. An Unknown Key Stops Startup before a Listener Opens. */
public final class Settings {

  private Settings() {
  }

  private static final ObjectMapper MAPPER = new ObjectMapper();

  /** A ValueRule Validates one declared Configuration Key. */
  public interface ValueRule extends Consumer<JsonNode> {
  }

  /** CheckIntegerRange Keeps JSON Numbers whole and Bounded. */
  public static ValueRule checkIntegerRange(int minimum, int maximum) {
    return value -> {
      if (value.isNull() || !value.isInt()) {
        throw new IllegalArgumentException("Must be an Integer");
      }
      if (value.asInt() < minimum || value.asInt() > maximum) {
        throw new IllegalArgumentException("Must be between " + minimum + " and " + maximum);
      }
    };
  }

  /** CheckTextValue Requires a nonempty String. */
  public static final ValueRule CHECK_TEXT_VALUE = value -> {
    if (!value.isTextual()) {
      throw new IllegalArgumentException("Must be a String");
    }
    if (value.asText().isBlank()) {
      throw new IllegalArgumentException("Must not be Empty");
    }
  };

  /** CheckChoiceValue Accepts only the Words this Runtime Implements. */
  public static ValueRule checkChoiceValue(String... allowed) {
    return value -> {
      CHECK_TEXT_VALUE.accept(value);
      for (String option : allowed) {
        if (option.equals(value.asText())) {
          return;
        }
      }
      throw new IllegalArgumentException("Must be one of " + String.join(", ", allowed));
    };
  }

  /** ReadDataFile Rejects Unknown Keys and Invalid Values before Merging. */
  public static Map<String, JsonNode> readDataFile(Path path, Map<String, ValueRule> rules) {
    JsonNode document;
    try {
      document = MAPPER.readTree(Files.readString(path));
    } catch (IOException failure) {
      throw new IllegalStateException(path + ": " + failure.getMessage(), failure);
    }

    if (document == null || !document.isObject()) {
      throw new IllegalStateException(path + " Must Hold an Object");
    }

    Map<String, JsonNode> values = new LinkedHashMap<>();
    document.fields().forEachRemaining(entry -> {
      ValueRule rule = rules.get(entry.getKey());
      if (rule == null) {
        throw new IllegalStateException(path + ": Unknown Key " + entry.getKey());
      }
      try {
        rule.accept(entry.getValue());
      } catch (IllegalArgumentException broken) {
        throw new IllegalStateException(path + ": " + entry.getKey() + " " + broken.getMessage(), broken);
      }
      values.put(entry.getKey(), entry.getValue());
    });

    return values;
  }

  /** DecodeDataValues Requires every declared Key before Returning typed Data. */
  public static <T> T decodeDataValues(Map<String, JsonNode> values, Map<String, ValueRule> rules, Class<T> shape) {
    for (String key : rules.keySet()) {
      if (!values.containsKey(key)) {
        throw new IllegalStateException(key + " is Missing");
      }
    }

    ObjectNode document = MAPPER.createObjectNode();
    values.forEach(document::set);

    return MAPPER.convertValue(document, shape);
  }

  /** ReadGlobalValues Loads Constants without an Environment Override Path. */
  public static <T> T readGlobalValues(Path path, Map<String, ValueRule> rules, Class<T> shape) {
    return decodeDataValues(readDataFile(path, rules), rules, shape);
  }

  /**
   * FindSchoolDataRoot Walks up until the Data Directories Appear.
   * A Binary Runs from java/ and a Test Runs from its own Directory.
   */
  public static Path findSchoolDataRoot() {
    Path root = Paths.get("").toAbsolutePath();

    for (int level = 0; level < 5 && root != null; level++, root = root.getParent()) {
      if (Files.exists(root.resolve("constants").resolve("school.json"))) {
        return root;
      }
    }

    throw new IllegalStateException("school Data Directories not Found above the Working Directory");
  }

  /** ReadProcessEnvironment Captures Startup Inputs once. */
  public static Map<String, String> readProcessEnvironment() {
    return new LinkedHashMap<>(System.getenv());
  }
}
