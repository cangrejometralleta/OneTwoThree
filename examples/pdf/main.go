package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// main Casts the Player, then Steps off the Stage.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "❌ Usage: go run . <source.md>")
		os.Exit(1)
	}

	out, err := buildPDF(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ %s\n", out)
}

// buildPDF Reads the Markdown, Shapes it, then Renders it beside the Source.
func buildPDF(source string) (string, error) {
	doc, err := ParseManifestoDocument(source)
	if err != nil {
		return "", err
	}

	out := outputPathFor(source)

	return out, RenderDocumentToPDF(doc, out)
}

// outputPathFor Swaps the Extension, so the PDF Lands beside its Source.
func outputPathFor(source string) string {
	return strings.TrimSuffix(source, filepath.Ext(source)) + ".pdf"
}
