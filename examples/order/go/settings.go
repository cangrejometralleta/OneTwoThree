package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// valueRule Validates one declared Configuration Key.
type valueRule func(json.RawMessage) error

// readDataFile Rejects Unknown Keys and Invalid Values before Merging.
func readDataFile(path string, rules map[string]valueRule) (map[string]json.RawMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var values map[string]json.RawMessage
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if values == nil {
		return nil, fmt.Errorf("%s Must Hold an Object", path)
	}
	for key, value := range values {
		rule, exists := rules[key]
		if !exists {
			return nil, fmt.Errorf("%s: Unknown Key %s", path, key)
		}
		if err := rule(value); err != nil {
			return nil, fmt.Errorf("%s: %s %w", path, key, err)
		}
	}

	return values, nil
}

// checkIntegerRange Keeps JSON Numbers whole and Bounded.
func checkIntegerRange(minimum, maximum int) valueRule {
	return func(raw json.RawMessage) error {
		var value int
		if string(raw) == "null" {
			return fmt.Errorf("Must not be Null")
		}
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("Must be an Integer")
		}
		if value < minimum || value > maximum {
			return fmt.Errorf("Must be between %d and %d", minimum, maximum)
		}
		return nil
	}
}

// checkTextValue Requires a nonempty String.
func checkTextValue(raw json.RawMessage) error {
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return fmt.Errorf("Must be a String")
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("Must not be Empty")
	}
	return nil
}

// decodeDataValues Requires every declared Key before Returning typed Data.
func decodeDataValues[T any](values map[string]json.RawMessage, rules map[string]valueRule) (T, error) {
	var result T
	for key := range rules {
		if _, exists := values[key]; !exists {
			return result, fmt.Errorf("%s is Missing", key)
		}
	}

	data, err := json.Marshal(values)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)

	return result, err
}

// readGlobalValues Loads Constants without an Environment Override Path.
func readGlobalValues[T any](path string, rules map[string]valueRule) T {
	values, err := readDataFile(path, rules)
	if err != nil {
		panic(err)
	}

	result, err := decodeDataValues[T](values, rules)
	if err != nil {
		panic(err)
	}

	return result
}

// readEnvironmentFiles Applies Defaults before the selected Environment.
func readEnvironmentFiles(root, concern, environment string, rules map[string]valueRule) (map[string]json.RawMessage, error) {
	if environment != "development" && environment != "production" {
		return nil, fmt.Errorf("APP_ENV Must be development or production")
	}

	values, err := readDataFile(filepath.Join(root, "config", concern+".defaults.json"), rules)
	if err != nil {
		return nil, err
	}
	overrides, err := readDataFile(filepath.Join(root, "config", concern+"."+environment+".json"), rules)
	if err != nil {
		return nil, err
	}
	for key, value := range overrides {
		values[key] = value
	}

	return values, nil
}

// checkEnvironmentKeys Refuses undeclared Variables in this Service Namespace.
func checkEnvironmentKeys(environment map[string]string, prefix string, names map[string]string) error {
	for name := range environment {
		if _, exists := names[name]; strings.HasPrefix(name, prefix) && !exists {
			return fmt.Errorf("Unknown Variable %s; Global Constants Cannot be Overridden", name)
		}
	}

	return nil
}

// applyEnvironmentValues Validates declared Variables before Replacing Values.
func applyEnvironmentValues(values map[string]json.RawMessage, environment map[string]string, names map[string]string, rules map[string]valueRule) error {
	for name, key := range names {
		value, exists := environment[name]
		if !exists {
			continue
		}
		raw, err := encodeEnvironmentValue(value, rules[key])
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		values[key] = raw
	}

	return nil
}

// encodeEnvironmentValue Accepts Text or a decimal Integer, never implicit Booleans.
func encodeEnvironmentValue(value string, rule valueRule) (json.RawMessage, error) {
	raw, _ := json.Marshal(value)
	if rule(raw) == nil {
		return raw, nil
	}

	if !regexp.MustCompile(`^[0-9]+$`).MatchString(value) {
		return nil, fmt.Errorf("Invalid Value")
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("Invalid Integer")
	}
	raw, _ = json.Marshal(number)

	return raw, rule(raw)
}

// readProcessEnvironment Captures Startup Inputs once.
func readProcessEnvironment() map[string]string {
	values := map[string]string{}
	for _, entry := range os.Environ() {
		name, value, _ := strings.Cut(entry, "=")
		values[name] = value
	}
	return values
}
