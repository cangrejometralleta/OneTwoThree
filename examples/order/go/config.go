package main

// OrderConfig Holds Deployment Values; Secrets Have no File Key.
type OrderConfig struct {
	Port         int    `json:"port"`
	DatabasePath string `json:"databasePath"`
}

// LoadOrderConfig Resolves Startup Inputs before any Store or Listener Opens.
func LoadOrderConfig(root string, environment map[string]string) (OrderConfig, error) {
	rules := map[string]valueRule{
		"port":         checkIntegerRange(1, 65535),
		"databasePath": checkTextValue,
	}
	names := map[string]string{
		"ORDER_PORT":          "port",
		"ORDER_DATABASE_PATH": "databasePath",
	}
	if err := checkEnvironmentKeys(environment, "ORDER_", names); err != nil {
		return OrderConfig{}, err
	}

	values, err := readEnvironmentFiles(root, "order", environment["APP_ENV"], rules)
	if err != nil {
		return OrderConfig{}, err
	}
	if err := applyEnvironmentValues(values, environment, names, rules); err != nil {
		return OrderConfig{}, err
	}
	config, err := decodeDataValues[OrderConfig](values, rules)
	if err != nil {
		return config, err
	}

	return config, nil
}
