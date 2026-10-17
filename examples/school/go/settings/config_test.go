package settings

import (
	"os"
	"path/filepath"
	"testing"
)

var adapters = []string{"chi", "gin", "stdlib"}

// buildConfigFixture Keeps broken Files away from the checked-in Examples.
func buildConfigFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, folder := range []string{"config", "constants"} {
		if err := os.Mkdir(filepath.Join(root, folder), 0o700); err != nil {
			t.Fatal(err)
		}
	}

	for _, name := range []string{"config/school.defaults.json", "config/school.development.json",
		"config/school.production.json", "config/registry.json", "constants/school.json"} {
		data, err := os.ReadFile(filepath.Join("..", "..", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	return root
}

func readyEnvironment() map[string]string {
	return map[string]string{"APP_ENV": "production", "TOKEN_SECRET": "test-secret"}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestTheShippedFilesLoadForBothEnvironments(t *testing.T) {
	root := buildConfigFixture(t)
	environment := readyEnvironment()

	production, err := LoadSchoolConfig(root, environment, adapters)
	if err != nil || production.DatabasePath != "school.production.db" {
		t.Fatalf("production=%+v err=%v", production, err)
	}

	environment["APP_ENV"] = "development"
	development, err := LoadSchoolConfig(root, environment, adapters)
	if err != nil || development.DatabasePath != "school.db" || development.Port != 8080 ||
		development.TokenLifeSeconds != 600 || development.ServerAdapter != "stdlib" {
		t.Fatalf("development=%+v err=%v", development, err)
	}
}

func TestDefaultsThenTheEnvironmentFileThenTheVariables(t *testing.T) {
	root := buildConfigFixture(t)
	writeFile(t, filepath.Join(root, "config", "school.production.json"), `{"databasePath":"production.db","port":9000}`)
	environment := readyEnvironment()

	config, err := LoadSchoolConfig(root, environment, adapters)
	if err != nil || config.Port != 9000 || config.DatabasePath != "production.db" {
		t.Fatalf("the Environment File Must Replace Defaults, config=%+v err=%v", config, err)
	}

	environment["SCHOOL_PORT"] = "9100"
	environment["SERVER"] = "gin"
	config, err = LoadSchoolConfig(root, environment, adapters)
	if err != nil || config.Port != 9100 || config.ServerAdapter != "gin" {
		t.Fatalf("a Variable Must Replace the Environment File, config=%+v err=%v", config, err)
	}
}

func TestTheRegistryResolvesFromTheDataRootAndMustExist(t *testing.T) {
	root := buildConfigFixture(t)
	environment := readyEnvironment()

	config, err := LoadSchoolConfig(root, environment, adapters)
	if err != nil || config.RegistryPath != filepath.Join(root, "config", "registry.json") {
		t.Fatalf("RegistryPath=%s err=%v", config.RegistryPath, err)
	}

	environment["SCHOOL_REGISTRY_PATH"] = filepath.Join(root, "config", "absent.json")
	if _, err := LoadSchoolConfig(root, environment, adapters); err == nil {
		t.Fatal("A Registry File that does not Exist Must Stop Startup")
	}

	environment["SCHOOL_REGISTRY_PATH"] = "config"
	if _, err := LoadSchoolConfig(root, environment, adapters); err == nil {
		t.Fatal("A Directory is not a Registry File")
	}
}

func TestConfigRejectsInvalidInputs(t *testing.T) {
	cases := map[string]map[string]string{
		"missing environment":   {"APP_ENV": ""},
		"unknown environment":   {"APP_ENV": "staging"},
		"path traversal":        {"APP_ENV": "../production"},
		"empty override":        {"SCHOOL_PORT": ""},
		"signed port":           {"SCHOOL_PORT": "+80"},
		"spaced port":           {"SCHOOL_PORT": " 80"},
		"fractional port":       {"SCHOOL_PORT": "80.0"},
		"port too high":         {"SCHOOL_PORT": "65536"},
		"port zero":             {"SCHOOL_PORT": "0"},
		"empty path":            {"SCHOOL_DATABASE_PATH": " "},
		"empty registry path":   {"SCHOOL_REGISTRY_PATH": ""},
		"global variable":       {"SCHOOL_MINIMUM_AGE_YEARS": "1"},
		"unknown variable":      {"SCHOOL_TYPO": "1"},
		"missing secret":        {"TOKEN_SECRET": ""},
		"unknown adapter":       {"SERVER": "spring"},
		"empty adapter":         {"SERVER": ""},
		"lifetime zero":         {"SCHOOL_TOKEN_LIFE_SECONDS": "0"},
		"lifetime beyond a day": {"SCHOOL_TOKEN_LIFE_SECONDS": "86401"},
	}

	for name, overrides := range cases {
		root := buildConfigFixture(t)
		environment := readyEnvironment()
		for key, value := range overrides {
			environment[key] = value
		}

		if _, err := LoadSchoolConfig(root, environment, adapters); err == nil {
			t.Errorf("%s: Invalid Configuration Must Fail", name)
		}
	}

	unsecret := readyEnvironment()
	delete(unsecret, "TOKEN_SECRET")
	if _, err := LoadSchoolConfig(buildConfigFixture(t), unsecret, adapters); err == nil {
		t.Error("an absent Secret Must Fail")
	}
}

func TestConfigRejectsInvalidFiles(t *testing.T) {
	cases := map[string]string{
		"global file key":    `{"databasePath":"test.db","minimumAgeYears":1}`,
		"unknown file key":   `{"databasePath":"test.db","typo":1}`,
		"null key":           `{"databasePath":null}`,
		"invalid type":       `{"databasePath":"test.db","port":"9000"}`,
		"fractional":         `{"databasePath":"test.db","port":9000.5}`,
		"invalid file range": `{"databasePath":"test.db","port":0}`,
		"invalid json":       `{`,
		"null document":      `null`,
		"array document":     `[]`,
		"empty file":         ``,
		"second value":       `{} {}`,
		"secret in file":     `{"databasePath":"test.db","tokenSecret":"forbidden"}`,
	}

	for name, content := range cases {
		root := buildConfigFixture(t)
		writeFile(t, filepath.Join(root, "config", "school.production.json"), content)

		if _, err := LoadSchoolConfig(root, readyEnvironment(), adapters); err == nil {
			t.Errorf("%s: Invalid Configuration Must Fail", name)
		}
	}

	root := buildConfigFixture(t)
	os.Remove(filepath.Join(root, "config", "school.production.json"))
	if _, err := LoadSchoolConfig(root, readyEnvironment(), adapters); err == nil {
		t.Error("a missing Environment File Must Fail")
	}
}

func TestAMissingKeyStopsStartup(t *testing.T) {
	root := buildConfigFixture(t)
	writeFile(t, filepath.Join(root, "config", "school.defaults.json"), `{"port":8080}`)
	writeFile(t, filepath.Join(root, "config", "school.production.json"), `{}`)

	if _, err := LoadSchoolConfig(root, readyEnvironment(), adapters); err == nil {
		t.Fatal("a missing Key Must Fail")
	}
}

func TestAnOverrideNeverHidesAnInvalidDefault(t *testing.T) {
	root := buildConfigFixture(t)
	writeFile(t, filepath.Join(root, "config", "school.defaults.json"), `{"port":null}`)
	environment := readyEnvironment()
	environment["SCHOOL_PORT"] = "9100"

	if _, err := LoadSchoolConfig(root, environment, adapters); err == nil {
		t.Fatal("An Override Must not Hide Invalid Defaults")
	}
}

func TestAVariableOutsideTheNamespaceIsIgnored(t *testing.T) {
	environment := readyEnvironment()
	environment["SCHOOLING"] = "anything"
	environment["PATH"] = "/usr/bin"

	if config, err := LoadSchoolConfig(buildConfigFixture(t), environment, adapters); err != nil || config.Port != 8080 {
		t.Fatalf("config=%+v err=%v", config, err)
	}
}

func TestTheGlobalConstantsHoldTheEnrolmentAgeAndRefuseBrokenData(t *testing.T) {
	root := buildConfigFixture(t)
	if constants, err := ReadGlobalConstants(root); err != nil || constants.MinimumAgeYears != 18 {
		t.Fatalf("constants=%+v err=%v", constants, err)
	}

	for _, broken := range []string{`{"minimumAgeYears":null}`, `{}`, `{"minimumAgeYears":0}`, `{"minimumAgeYears":18,"port":1}`} {
		writeFile(t, filepath.Join(root, "constants", "school.json"), broken)
		if _, err := ReadGlobalConstants(root); err == nil {
			t.Errorf("%s Must Fail", broken)
		}
	}
}

func TestTheDataRootIsFoundAboveTheWorkingDirectory(t *testing.T) {
	root, err := FindSchoolDataRoot(".")
	want, _ := filepath.Abs(filepath.Join("..", ".."))

	if err != nil || root != want {
		t.Fatalf("root=%s want=%s err=%v", root, want, err)
	}
}
