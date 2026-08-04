package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOutputPathForSwapsTheExtension(t *testing.T) {
	cases := map[string]string{
		"source.md":          "source.pdf",
		"nested/dir/file.md": "nested/dir/file.pdf",
		"noextension":        "noextension.pdf",
		"two.dots.md":        "two.dots.pdf",
	}

	for in, want := range cases {
		if got := outputPathFor(in); got != want {
			t.Errorf("outputPathFor(%q): wanted %q, got %q", in, want, got)
		}
	}
}

func TestBuildPDFWritesBesideTheSource(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "sample.md")

	original, err := os.ReadFile("testdata/sample.md")
	if err != nil {
		t.Fatalf("the Sample Must Read: %v", err)
	}
	if err := os.WriteFile(source, original, 0o644); err != nil {
		t.Fatalf("the Copy Must Write: %v", err)
	}

	out, err := buildPDF(source)
	if err != nil {
		t.Fatalf("buildPDF Must Succeed, got %v", err)
	}

	if want := filepath.Join(dir, "sample.pdf"); out != want {
		t.Errorf("wanted the Output at %q, got %q", want, out)
	}
	assertIsPDF(t, out)
}

func TestBuildPDFRefusesAMissingSource(t *testing.T) {
	if _, err := buildPDF(filepath.Join(t.TempDir(), "missing.md")); err == nil {
		t.Fatal("a missing Source Must Refuse to Build")
	}
}
