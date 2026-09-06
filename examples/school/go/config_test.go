package main

import (
	"os"
	"path/filepath"
	"testing"
)

// buildConfigFixture Keeps broken Files away from the checked-in Examples.
func buildConfigFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "config"), 0700); err != nil {
		t.Fatal(err)
	}

	for _, layer := range []string{"defaults", "development", "production"} {
		name := "school." + layer + ".json"
		data, err := os.ReadFile(filepath.Join("..", "config", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "config", name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}

	return root
}

func TestConfigPrecedence(t *testing.T) {
	root := buildConfigFixture(t)
	path := filepath.Join(root, "config", "school.production.json")
	if err := os.WriteFile(path, []byte(`{"databasePath":"production.db","port":9000}`), 0600); err != nil {
		t.Fatal(err)
	}
	environment := map[string]string{"APP_ENV": "production", "TOKEN_SECRET": "test-secret"}

	config, err := LoadSchoolConfig(root, environment)
	if err != nil {
		t.Fatal(err)
	}
	if config.Port != 9000 || config.DatabasePath != "production.db" {
		t.Fatal("Environment File Must Replace Defaults")
	}
	environment["SCHOOL_PORT"] = "9100"
	config, err = LoadSchoolConfig(root, environment)
	if err != nil {
		t.Fatal(err)
	}

	if config.Port != 9100 {
		t.Fatal("Variable Must Replace Environment File")
	}
	if ReadSchoolConstants().MinimumAgeYears != 18 {
		t.Fatal("Deployment Must Preserve Global Constants")
	}
}

func TestConfigRejectsInvalidInputs(t *testing.T) {
	cases := []struct {
		name        string
		environment map[string]string
		file        string
		missing     bool
	}{
		{name: "missing environment", environment: map[string]string{"APP_ENV": ""}},
		{name: "unknown environment", environment: map[string]string{"APP_ENV": "staging"}},
		{name: "path traversal", environment: map[string]string{"APP_ENV": "../production"}},
		{name: "empty override", environment: map[string]string{"SCHOOL_PORT": ""}},
		{name: "invalid port", environment: map[string]string{"SCHOOL_PORT": "65536"}},
		{name: "fractional port", environment: map[string]string{"SCHOOL_PORT": "3.5"}},
		{name: "empty path", environment: map[string]string{"SCHOOL_DATABASE_PATH": " "}},
		{name: "global variable", environment: map[string]string{"SCHOOL_MINIMUM_AGE_YEARS": "1"}},
		{name: "unknown variable", environment: map[string]string{"SCHOOL_TYPO": "1"}},
		{name: "global file key", file: `{"databasePath":"test.db","minimumAgeYears":1}`},
		{name: "unknown file key", file: `{"databasePath":"test.db","typo":1}`},
		{name: "null key", file: `{"databasePath":null}`},
		{name: "invalid type", file: `{"databasePath":"test.db","port":"9000"}`},
		{name: "invalid file range", file: `{"databasePath":"test.db","port":0}`},
		{name: "missing key", file: `{}`},
		{name: "invalid json", file: `{`},
		{name: "null document", file: `null`},
		{name: "array document", file: `[]`},
		{name: "missing file", missing: true},
		{name: "missing secret", environment: map[string]string{"TOKEN_SECRET": ""}},
		{name: "unknown adapter", environment: map[string]string{"SERVER": "missing"}},
		{name: "invalid lifetime", environment: map[string]string{"SCHOOL_TOKEN_LIFE_SECONDS": "0"}},
		{name: "secret in file", file: `{"databasePath":"test.db","tokenSecret":"forbidden"}`},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root := buildConfigFixture(t)
			path := filepath.Join(root, "config", "school.production.json")
			if test.file != "" {
				if err := os.WriteFile(path, []byte(test.file), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if test.missing {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			}
			environment := map[string]string{"APP_ENV": "production", "TOKEN_SECRET": "test-secret"}
			for key, value := range test.environment {
				environment[key] = value
			}

			if _, err := LoadSchoolConfig(root, environment); err == nil {
				t.Fatal("Invalid Configuration Must Fail")
			}
		})
	}
}

func TestConfigValidatesBeforeOverrides(t *testing.T) {
	root := buildConfigFixture(t)
	path := filepath.Join(root, "config", "school.defaults.json")
	if err := os.WriteFile(path, []byte(`{"port":null}`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadSchoolConfig(root, map[string]string{"APP_ENV": "production", "SCHOOL_PORT": "9100", "TOKEN_SECRET": "test-secret"})

	if err == nil {
		t.Fatal("An Override Must not Hide Invalid Defaults")
	}
}

func TestGlobalConstantsRejectInvalidData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "constants.json")
	if err := os.WriteFile(path, []byte(`{"minimumAgeYears":null}`), 0600); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("Invalid Constants Must Fail at Startup")
		}
	}()

	readGlobalValues[map[string]int](path, map[string]valueRule{"minimumAgeYears": checkIntegerRange(1, 200)})
}
