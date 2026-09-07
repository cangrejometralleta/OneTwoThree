package com.example.school;

import com.example.settings.Settings;

import java.util.Map;

/**
 * SchoolConstants Holds shared Meaning, independent of Deployment.
 * Keep the Snapshot private; Readers Receive a Value Copy.
 */
public final class Constants {

  private Constants() {
  }

  /** OldestPlausibleAge Bounds the Constant File, not the Business. */
  private static final int OLDEST_PLAUSIBLE_AGE = 150;

  public record SchoolConstants(int minimumAgeYears) {
  }

  private static final SchoolConstants SNAPSHOT = Settings.readGlobalValues(
      Settings.findSchoolDataRoot().resolve("constants").resolve("school.json"),
      Map.of("minimumAgeYears", Settings.checkIntegerRange(1, OLDEST_PLAUSIBLE_AGE)),
      SchoolConstants.class);

  /** ReadSchoolConstants Returns the Startup Snapshot without a Mutation Path. */
  public static SchoolConstants readSchoolConstants() {
    return SNAPSHOT;
  }
}
