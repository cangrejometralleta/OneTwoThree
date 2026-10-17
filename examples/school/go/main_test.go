package main

import (
	"os/exec"
	"strings"
	"testing"
)

const module = "github.com/cangrejometralleta/OneTwoThree/examples/school/go/"

// The Compiler Holds the Boundary, and go list Reads it: the Core Names no Vendor,
// no Socket and no Wire Package.
func TestTheCoreNamesNoVendorAndNoSocket(t *testing.T) {
	output, err := exec.Command("go", "list", "-deps", "./school").Output()
	if err != nil {
		t.Fatal(err)
	}

	allowed := map[string]bool{module + "school": true, module + "faults": true}
	for _, dependency := range strings.Fields(string(output)) {
		vendor := strings.Contains(dependency, ".") && !allowed[dependency]
		if vendor || strings.HasPrefix(dependency, "net") {
			t.Errorf("the Core Reaches %s", dependency)
		}
	}
}

// Each Vendor Lives in one Package: count the Packages that Import it.
func TestEveryVendorLivesInOnePackage(t *testing.T) {
	homes := map[string]string{
		"gorm.io/gorm":             module + "store",
		"github.com/go-chi/chi/v5": module + "serving/chiserver",
		"github.com/gin-gonic/gin": module + "serving/ginserver",
		"gorm.io/driver/sqlite":    module + "store",
	}

	output, err := exec.Command("go", "list", "-f", `{{.ImportPath}} {{join .Imports " "}}`, "./...").Output()
	if err != nil {
		t.Fatal(err)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		fields := strings.Fields(line)
		for _, imported := range fields[1:] {
			if home, vendor := homes[imported]; vendor && fields[0] != home {
				t.Errorf("%s Imports %s, which Lives in %s", fields[0], imported, home)
			}
		}
	}
}
