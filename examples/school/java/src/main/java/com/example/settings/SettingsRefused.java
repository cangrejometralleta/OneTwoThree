package com.example.settings;

/** SettingsRefused Stops Startup and Says what Broke. */
public final class SettingsRefused extends Exception {

  public SettingsRefused(String reason) {
    super(reason);
  }
}
