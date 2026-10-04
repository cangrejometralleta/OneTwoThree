package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The Bounds a Deployment Value Must Fall inside.
// Read alone, 65535 Says nothing and HighestPort Says the Rule.
const (
	LowestPort    = 1
	HighestPort   = 65535
	SecondsInADay = 60 * 60 * 24
)

// SchoolConfig Holds Deployment Values; Secrets Have no File Key.
type SchoolConfig struct {
	Port             int    `json:"port"`
	DatabasePath     string `json:"databasePath"`
	RegistryPath     string `json:"registryPath"`
	TokenLifeSeconds int    `json:"tokenLifeSeconds"`
	ServerAdapter    string `json:"serverAdapter"`
	TokenSecret      string `json:"-"`
}

// LoadSchoolConfig Resolves Startup Inputs before any Store or Listener Opens.
func LoadSchoolConfig(root string, environment map[string]string) (SchoolConfig, error) {
	rules := map[string]ValueRule{
		"port":             CheckIntegerRange(LowestPort, HighestPort),
		"databasePath":     CheckTextValue,
		"registryPath":     CheckTextValue,
		"tokenLifeSeconds": CheckIntegerRange(1, SecondsInADay),
		"serverAdapter":    checkServerAdapter,
	}
	names := map[string]string{
		"SCHOOL_PORT":               "port",
		"SCHOOL_DATABASE_PATH":      "databasePath",
		"SCHOOL_REGISTRY_PATH":      "registryPath",
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

	config.RegistryPath, err = findRegistryFile(root, config.RegistryPath)
	if err != nil {
		return config, err
	}

	config.TokenSecret = environment["TOKEN_SECRET"]
	if strings.TrimSpace(config.TokenSecret) == "" {
		return config, fmt.Errorf("TOKEN_SECRET is Missing")
	}

	return config, nil
}

// findRegistryFile Resolves a relative Path from the Data Root, where the
// Binary and a Test Agree, and Refuses a File that does not Exist.
func findRegistryFile(root, path string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}

	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return "", fmt.Errorf("registryPath %s Must be a File that Exists", path)
	}

	return path, nil
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
