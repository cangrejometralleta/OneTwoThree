package com.example.settings;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;
import java.util.Set;

/**
 * SchoolConfig Holds the Deployment Values, Validated once at Startup.
 * The Secret has no File Key, and never Prints.
 */
public record SchoolConfig(int port, Path databasePath, Path registryPath, long tokenLifeSeconds,
    String serverAdapter, String tokenSecret) {

  /** The Bounds a Deployment Value Must Fall inside. Read alone, 65535 Says nothing. */
  public static final int LOWEST_PORT = 1;
  public static final int HIGHEST_PORT = 65535;
  public static final long SECONDS_IN_A_DAY = 60L * 60 * 24;

  private static final String NAMESPACE = "SCHOOL_";

  /**
   * loadSchoolConfig Resolves Defaults, then the Environment File, then the Variables.
   * The Adapters a Runtime Implements Arrive from the Composition; this Loader Names none.
   */
  public static SchoolConfig loadSchoolConfig(Path root, Map<String, String> environment, Set<String> adapters)
      throws SettingsRefused {
    Map<String, ValueRule> rules = Map.of(
        "port", ValueRule.checkIntegerRange(LOWEST_PORT, HIGHEST_PORT),
        "databasePath", ValueRule.checkTextValue(),
        "registryPath", ValueRule.checkTextValue(),
        "tokenLifeSeconds", ValueRule.checkIntegerRange(1, SECONDS_IN_A_DAY),
        "serverAdapter", ValueRule.checkOneOf(adapters));
    Map<String, String> names = Map.of(
        "SCHOOL_PORT", "port",
        "SCHOOL_DATABASE_PATH", "databasePath",
        "SCHOOL_REGISTRY_PATH", "registryPath",
        "SCHOOL_TOKEN_LIFE_SECONDS", "tokenLifeSeconds",
        "SERVER", "serverAdapter");

    Settings.checkEnvironmentKeys(environment, NAMESPACE, names);
    Map<String, Object> values = Settings.readEnvironmentFiles(root, "school", environment.get("APP_ENV"), rules);
    Settings.applyEnvironmentValues(values, environment, names, rules);
    Settings.requireEveryKey(values, rules);

    String secret = environment.getOrDefault("TOKEN_SECRET", "");
    if (secret.isBlank()) {
      throw new SettingsRefused("TOKEN_SECRET is Missing");
    }

    return new SchoolConfig(((BigInteger) values.get("port")).intValueExact(),
        Path.of((String) values.get("databasePath")),
        findRegistryFile(root, (String) values.get("registryPath")),
        ((BigInteger) values.get("tokenLifeSeconds")).longValueExact(),
        (String) values.get("serverAdapter"), secret);
  }

  /** findRegistryFile Resolves a relative Path from the Data Root, and Refuses a File that does not Exist. */
  private static Path findRegistryFile(Path root, String path) throws SettingsRefused {
    Path resolved = root.resolve(path);
    if (!Files.isRegularFile(resolved)) {
      throw new SettingsRefused("registryPath " + resolved + " Must be a File that Exists");
    }

    return resolved;
  }

  @Override
  public String toString() {
    return "SchoolConfig[port=" + port + ", databasePath=" + databasePath + ", registryPath=" + registryPath
        + ", tokenLifeSeconds=" + tokenLifeSeconds + ", serverAdapter=" + serverAdapter + ", tokenSecret=***]";
  }
}
