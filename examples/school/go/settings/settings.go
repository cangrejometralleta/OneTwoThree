package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// dataRootSearchLevels Bounds the Walk up from the working Directory.
const dataRootSearchLevels = 5

// ValueRule Validates one declared Configuration Key.
type ValueRule func(json.RawMessage) error

// readDataFile Rejects Unknown Keys and Invalid Values before Merging.
func readDataFile(path string, rules map[string]ValueRule) (map[string]json.RawMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var values map[string]json.RawMessage
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&values); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%s Holds a second Value", path)
	}
	if values == nil {
		return nil, fmt.Errorf("%s Must Hold an Object", path)
	}

	for key, value := range values {
		rule, declared := rules[key]
		if !declared {
			return nil, fmt.Errorf("%s: Unknown Key %s", path, key)
		}
		if err := rule(value); err != nil {
			return nil, fmt.Errorf("%s: %s %w", path, key, err)
		}
	}

	return values, nil
}

var jsonWhole = regexp.MustCompile(`^-?(0|[1-9][0-9]*)$`)

// CheckIntegerRange Keeps JSON Numbers whole and Bounded.
func CheckIntegerRange(minimum, maximum int64) ValueRule {
	return func(raw json.RawMessage) error {
		var value int64
		if !jsonWhole.Match(raw) || json.Unmarshal(raw, &value) != nil {
			return fmt.Errorf("Must be an Integer")
		}
		if value < minimum || value > maximum {
			return fmt.Errorf("Must be between %d and %d", minimum, maximum)
		}

		return nil
	}
}

// CheckTextValue Requires a nonempty String.
func CheckTextValue(raw json.RawMessage) error {
	var value string
	if len(raw) == 0 || raw[0] != '"' || json.Unmarshal(raw, &value) != nil {
		return fmt.Errorf("Must be a String")
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("Must not be Empty")
	}

	return nil
}

// CheckOneOf Accepts only the Names it is Given.
func CheckOneOf(names []string) ValueRule {
	return func(raw json.RawMessage) error {
		var value string
		if len(raw) > 0 && raw[0] == '"' && json.Unmarshal(raw, &value) == nil {
			for _, name := range names {
				if value == name {
					return nil
				}
			}
		}

		return fmt.Errorf("Must be one of %v", names)
	}
}

// requireEveryKey Refuses a Set of Values that Lacks a declared Key.
func requireEveryKey(values map[string]json.RawMessage, rules map[string]ValueRule) error {
	for key := range rules {
		if _, present := values[key]; !present {
			return fmt.Errorf("%s is Missing", key)
		}
	}

	return nil
}

// GlobalConstants Holds the Values no Deployment can Override, Read from constants/.
type GlobalConstants struct {
	MinimumAgeYears int64 `json:"minimumAgeYears"`
}

// ReadGlobalConstants Loads Constants without an Override Path; no Variable Reaches them.
func ReadGlobalConstants(root string) (GlobalConstants, error) {
	rules := map[string]ValueRule{"minimumAgeYears": CheckIntegerRange(1, 200)}
	values, err := readDataFile(filepath.Join(root, "constants", "school.json"), rules)
	if err != nil {
		return GlobalConstants{}, err
	}
	if err := requireEveryKey(values, rules); err != nil {
		return GlobalConstants{}, err
	}

	var constants GlobalConstants

	return constants, json.Unmarshal(values["minimumAgeYears"], &constants.MinimumAgeYears)
}

// readEnvironmentFiles Applies the Defaults, then the selected Environment. Both Files must Exist.
func readEnvironmentFiles(root, concern, environment string, rules map[string]ValueRule) (map[string]json.RawMessage, error) {
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

// checkEnvironmentKeys Refuses undeclared Variables in this Service's Namespace.
func checkEnvironmentKeys(environment map[string]string, prefix string, names map[string]string) error {
	for name := range environment {
		if _, declared := names[name]; strings.HasPrefix(name, prefix) && !declared {
			return fmt.Errorf("Unknown Variable %s; Global Constants Cannot be Overridden", name)
		}
	}

	return nil
}

// applyEnvironmentValues Validates declared Variables before Replacing Values. Set and Empty is Set.
func applyEnvironmentValues(values map[string]json.RawMessage, environment, names map[string]string, rules map[string]ValueRule) error {
	for name, key := range names {
		value, present := environment[name]
		if !present {
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

var digitsOnly = regexp.MustCompile(`^[0-9]+$`)

// encodeEnvironmentValue Accepts Text or Digits only, so +80, " 80" and 80.0 Stop Startup.
func encodeEnvironmentValue(value string, rule ValueRule) (json.RawMessage, error) {
	text, _ := json.Marshal(value)
	if rule(text) == nil {
		return text, nil
	}

	if !digitsOnly.MatchString(value) {
		return nil, fmt.Errorf("Invalid Value")
	}

	raw := json.RawMessage(strings.TrimLeft(value, "0"))
	if len(raw) == 0 {
		raw = json.RawMessage("0")
	}

	return raw, rule(raw)
}

// FindSchoolDataRoot Walks up from a Directory until constants/ Appears.
// The Binary Runs from go/, and a Test from its own Package.
func FindSchoolDataRoot(start string) (string, error) {
	root, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for level := 0; level < dataRootSearchLevels; level++ {
		if info, err := os.Stat(filepath.Join(root, "constants", "school.json")); err == nil && !info.IsDir() {
			return root, nil
		}
		root = filepath.Dir(root)
	}

	return "", fmt.Errorf("school Data Directories not Found above %s", start)
}

// ReadProcessEnvironment Captures Startup Inputs once.
func ReadProcessEnvironment() map[string]string {
	values := map[string]string{}
	for _, entry := range os.Environ() {
		name, value, _ := strings.Cut(entry, "=")
		values[name] = value
	}

	return values
}
