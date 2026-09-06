package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SchoolConfig Holds Deployment Values; Secrets Have no File Key.
type SchoolConfig struct {
	Port             int    `json:"port"`
	DatabasePath     string `json:"databasePath"`
	TokenLifeSeconds int    `json:"tokenLifeSeconds"`
	ServerAdapter    string `json:"serverAdapter"`
	TokenSecret      string `json:"-"`
}

// LoadSchoolConfig Resolves Startup Inputs before any Store or Listener Opens.
func LoadSchoolConfig(root string, environment map[string]string) (SchoolConfig, error) {
	rules := map[string]valueRule{
		"port":             checkIntegerRange(1, 65535),
		"databasePath":     checkTextValue,
		"tokenLifeSeconds": checkIntegerRange(1, 86400),
		"serverAdapter":    checkServerAdapter,
	}
	names := map[string]string{
		"SCHOOL_PORT":               "port",
		"SCHOOL_DATABASE_PATH":      "databasePath",
		"SCHOOL_TOKEN_LIFE_SECONDS": "tokenLifeSeconds",
		"SERVER":                    "serverAdapter",
	}
	if err := checkEnvironmentKeys(environment, "SCHOOL_", names); err != nil {
		return SchoolConfig{}, err
	}

	values, err := readEnvironmentFiles(root, "school", environment["APP_ENV"], rules)
	if err != nil {
		return SchoolConfig{}, err
	}
	if err := applyEnvironmentValues(values, environment, names, rules); err != nil {
		return SchoolConfig{}, err
	}
	config, err := decodeDataValues[SchoolConfig](values, rules)
	if err != nil {
		return config, err
	}

	config.TokenSecret = environment["TOKEN_SECRET"]
	if strings.TrimSpace(config.TokenSecret) == "" {
		return config, fmt.Errorf("TOKEN_SECRET is Missing")
	}

	return config, nil
}

// checkServerAdapter Accepts only Adapters this Runtime Implements.
func checkServerAdapter(raw json.RawMessage) error {
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return fmt.Errorf("Must be a String")
	}
	switch value {
	case "stdlib", "chi", "gin":
		return nil
	default:
		return fmt.Errorf("Must be stdlib, chi or gin")
	}
}
