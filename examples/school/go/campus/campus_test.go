package campus

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
)

// One Office Fills both Ports, and the Compiler Checks it.
var (
	_ app.RutRegistry        = Office{}
	_ app.EnrollmentNotifier = Office{}
)

func writeRegistry(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "registry.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestConfirmRutKnowsOnlyWhatTheFileHoldsInLowerCase(t *testing.T) {
	office, err := OpenOffice(writeRegistry(t, `["11111111-1","10000013-K","11111111-1"]`), nil)
	if err != nil {
		t.Fatal(err)
	}

	known, _ := office.ConfirmRut("11111111-1")
	upper, _ := office.ConfirmRut("10000013-k")
	stranger, _ := office.ConfirmRut("22222222-2")

	if !known || !upper || stranger {
		t.Fatalf("known=%v upper=%v stranger=%v", known, upper, stranger)
	}
}

func TestOpenOfficeStopsOnWhatItCannotTrust(t *testing.T) {
	cases := map[string]string{
		"missing file":   filepath.Join(t.TempDir(), "absent.json"),
		"empty":          writeRegistry(t, ``),
		"broken JSON":    writeRegistry(t, `["11111111-1"`),
		"not an array":   writeRegistry(t, `{"ruts": []}`),
		"null":           writeRegistry(t, `null`),
		"a second value": writeRegistry(t, `[] []`),
		"a bad check":    writeRegistry(t, `["11111111-9"]`),
		"a non RUT text": writeRegistry(t, `["hello"]`),
		"not a text":     writeRegistry(t, `[11111111]`),
	}

	for name, path := range cases {
		if _, err := OpenOffice(path, nil); err == nil {
			t.Errorf("%s: wanted a startup Failure", name)
		}
	}
}

func TestAnEmptyRegistryIsAllowed(t *testing.T) {
	if _, err := OpenOffice(writeRegistry(t, `[]`), nil); err != nil {
		t.Fatal(err)
	}
}

func TestTheShippedRegistryHoldsOnlyRealRuts(t *testing.T) {
	office, err := OpenOffice(filepath.Join("..", "..", "config", "registry.json"), nil)
	if err != nil {
		t.Fatal(err)
	}

	if known, _ := office.ConfirmRut("12345678-5"); !known {
		t.Fatal("the shipped Registry must Know its example RUT")
	}
}

func TestAnAnnouncementNamesTheStudentByIdOnly(t *testing.T) {
	var out bytes.Buffer
	office := Office{Logger: slog.New(slog.NewTextHandler(&out, nil))}

	err := office.AnnounceEnrollment(school.Student{ID: 7, Rut: "11111111-1", Name: "Ana", Course: 1})

	if err != nil || !strings.Contains(out.String(), `msg="student 7 is enrolled"`) {
		t.Fatalf("err=%v log=%s", err, out.String())
	}
	if strings.Contains(out.String(), "11111111-1") || strings.Contains(out.String(), "Ana") {
		t.Fatalf("a Notice must not Leak the RUT or the Name: %s", out.String())
	}
}
