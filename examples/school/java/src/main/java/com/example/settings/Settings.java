package com.example.settings;

import com.example.codec.Json;

import java.io.IOException;
import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Pattern;

/**
 * Settings is the strict JSON Loader: Unknown Keys, nulls, wrong Types and
 * missing required Values Stop Startup before anything Opens.
 */
public final class Settings {

  /** DATA_ROOT_SEARCH_LEVELS Bounds the Walk up from the working Directory. */
  private static final int DATA_ROOT_SEARCH_LEVELS = 5;

  private static final Pattern WHOLE_NUMBER = Pattern.compile("^[0-9]+$");

  private Settings() {
  }

  /** readDataFile Rejects Unknown Keys and Invalid Values before Merging. */
  static Map<String, Object> readDataFile(Path path, Map<String, ValueRule> rules) throws SettingsRefused {
    byte[] data;
    try {
      data = Files.readAllBytes(path);
    } catch (IOException missing) {
      throw new SettingsRefused(path + " cannot be Read");
    }

    Object read = Json.readJsonText(data).orElseThrow(() -> new SettingsRefused(path + " is not one JSON value"));
    if (!(read instanceof Map<?, ?> values)) {
      throw new SettingsRefused(path + " Must Hold an Object");
    }

    Map<String, Object> checked = new LinkedHashMap<>();
    for (Map.Entry<?, ?> entry : values.entrySet()) {
      String key = (String) entry.getKey();
      ValueRule rule = rules.get(key);
      if (rule == null) {
        throw new SettingsRefused(path + ": Unknown Key " + key);
      }
      checkValue(rule, entry.getValue(), path + ": " + key);
      checked.put(key, entry.getValue());
    }

    return checked;
  }

  /** requireEveryKey Refuses a Set of Values that Lacks a declared Key. */
  static void requireEveryKey(Map<String, Object> values, Map<String, ValueRule> rules) throws SettingsRefused {
    for (String key : rules.keySet()) {
      if (!values.containsKey(key)) {
        throw new SettingsRefused(key + " is Missing");
      }
    }
  }

  /** readGlobalConstants Loads Constants without an Override Path; no Variable Reaches them. */
  public static GlobalConstants readGlobalConstants(Path root) throws SettingsRefused {
    Map<String, ValueRule> rules = Map.of("minimumAgeYears", ValueRule.checkIntegerRange(1, 200));
    Map<String, Object> values = readDataFile(root.resolve("constants").resolve("school.json"), rules);
    requireEveryKey(values, rules);

    return new GlobalConstants(((BigInteger) values.get("minimumAgeYears")).intValueExact());
  }

  /** readEnvironmentFiles Applies the Defaults, then the selected Environment. Both Files must Exist. */
  static Map<String, Object> readEnvironmentFiles(Path root, String concern, String environment,
      Map<String, ValueRule> rules) throws SettingsRefused {
    if (!"development".equals(environment) && !"production".equals(environment)) {
      throw new SettingsRefused("APP_ENV Must be development or production");
    }

    Path config = root.resolve("config");
    Map<String, Object> values = readDataFile(config.resolve(concern + ".defaults.json"), rules);
    values.putAll(readDataFile(config.resolve(concern + "." + environment + ".json"), rules));

    return values;
  }

  /** checkEnvironmentKeys Refuses undeclared Variables in this Service's Namespace. */
  static void checkEnvironmentKeys(Map<String, String> environment, String prefix, Map<String, String> names)
      throws SettingsRefused {
    for (String name : environment.keySet()) {
      if (name.startsWith(prefix) && !names.containsKey(name)) {
        throw new SettingsRefused("Unknown Variable " + name + "; Global Constants Cannot be Overridden");
      }
    }
  }

  /** applyEnvironmentValues Validates declared Variables before Replacing Values. Set and Empty is Set. */
  static void applyEnvironmentValues(Map<String, Object> values, Map<String, String> environment,
      Map<String, String> names, Map<String, ValueRule> rules) throws SettingsRefused {
    for (Map.Entry<String, String> declared : names.entrySet()) {
      String text = environment.get(declared.getKey());
      if (text != null) {
        values.put(declared.getValue(), encodeEnvironmentValue(text, rules.get(declared.getValue()),
            declared.getKey()));
      }
    }
  }

  /** encodeEnvironmentValue Accepts Text or Digits only, so +80, " 80" and 80.0 Stop Startup. */
  private static Object encodeEnvironmentValue(String text, ValueRule rule, String name) throws SettingsRefused {
    try {
      rule.check(text);

      return text;
    } catch (SettingsRefused notText) {
      if (!WHOLE_NUMBER.matcher(text).matches()) {
        throw new SettingsRefused(name + ": " + notText.getMessage());
      }
    }

    BigInteger number = new BigInteger(text);
    checkValue(rule, number, name);

    return number;
  }

  private static void checkValue(ValueRule rule, Object value, String where) throws SettingsRefused {
    try {
      rule.check(value);
    } catch (SettingsRefused refused) {
      throw new SettingsRefused(where + " " + refused.getMessage());
    }
  }

  /** findSchoolDataRoot Walks up until constants/ Appears: the Jar Runs from java/, a Test from java/ too. */
  public static Path findSchoolDataRoot(Path start) throws SettingsRefused {
    Path root = start.toAbsolutePath();
    for (int level = 0; level < DATA_ROOT_SEARCH_LEVELS && root != null; level++) {
      if (Files.isRegularFile(root.resolve("constants").resolve("school.json"))) {
        return root;
      }
      root = root.getParent();
    }

    throw new SettingsRefused("school Data Directories not Found above " + start.toAbsolutePath());
  }
}
