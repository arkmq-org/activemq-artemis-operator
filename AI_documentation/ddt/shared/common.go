package shared

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

// Documentation files to process
var DocFiles = []string{
	"operator_architecture.md",
	"operator_conventions.md",
	"operator_defaults_reference.md",
	"contribution_guide.md",
	"tdd_index.md",
}

// GetDocDir returns the documentation directory
// Assumes tool is run from AI_documentation/ddt/ or AI_documentation/
func GetDocDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if we're in the ddt subdirectory
	if _, err := os.Stat("../glossary_terms.json"); err == nil {
		// We're in ddt/ directory, parent is AI_documentation
		return "..", nil
	}

	// Check if we're in AI_documentation
	if _, err := os.Stat("glossary_terms.json"); err == nil {
		return ".", nil
	}

	return cwd, nil
}

// FindDocFiles finds all markdown files in the documentation directory
func FindDocFiles(docDir string) ([]string, error) {
	var existing []string
	for _, file := range DocFiles {
		fullPath := fmt.Sprintf("%s/%s", docDir, file)
		if _, err := os.Stat(fullPath); err == nil {
			existing = append(existing, fullPath)
		}
	}

	if len(existing) == 0 {
		return nil, fmt.Errorf("no documentation files found in %s", docDir)
	}

	return existing, nil
}

// SetupPterm configures pterm based on quiet/no-style flags
func SetupPterm(quiet, noStyle bool) {
	if quiet || noStyle {
		pterm.DisableStyling()
		pterm.DisableColor()
	}
}

// FormatNumber formats a number with thousand separators
func FormatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	var result []rune
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, c)
	}
	return string(result)
}
