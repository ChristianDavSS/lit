package tests

import (
	"CLI_App/cmd/adapters/analysis/languages"
	"CLI_App/cmd/domain"
	"fmt"
	"testing"
)

// TestScanning runs the scanner with mocks to see if it runs properly or if it breaks.
func TestScanning(t *testing.T) {
	mocks := []struct {
		analyzer languages.FileAnalyzer
		code     []byte
		filePath string
	}{
		{
			*languages.NewFileAnalyzer(false, domain.SnakeCase),
			[]byte("hola_mundo_soy_yo = \"hola\"; hola_mundo_soy_yo+=\"\"\n\nres = hola_mundo_soy_yo.split(\"a\")\n\nfor i in hola_mundo_soy_yo:\n    hola_mundo_soy_yo += \"a\";"),
			"test.py",
		},
		{
			*languages.NewFileAnalyzer(false, domain.SnakeCase),
			[]byte("hola_mundo_soy_yo = \"hola\"; hola_mundo_soy_yo+=\"\"\n\nres = hola_mundo_soy_yo.split(\"a\")\n\nfor i in hola_mundo_soy_yo:\n    hola_mundo_soy_yo += \"a\";"),
			"test.java",
		},
	}

	for i, m := range mocks {
		r := m.analyzer.AnalyzeFile(m.filePath, m.code)
		if r == nil {
			t.Errorf("- > Mock %d FAILED. The file extension isn't supported yet.\n", i+1)
			continue
		}
		fmt.Printf("- > Mock %d executed successfully. Returned %d functions\n", i+1, len(r))
	}
}
