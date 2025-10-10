package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/arkmq-org/ddt/shared"
)

// GlossaryTerm represents a single term in the glossary
type GlossaryTerm struct {
	Name           string   `json:"name"`
	Definition     string   `json:"definition"`
	SearchPatterns []string `json:"searchPatterns"`
	Category       string   `json:"category"`
}

// TermsConfig represents the JSON config structure
type TermsConfig struct {
	Terms []GlossaryTerm `json:"terms"`
}

// SectionMatch represents a matched section in a document
type SectionMatch struct {
	Anchor string
	Title  string
}

// Category order for output
var categoryOrder = []string{
	"Operator Core Concepts",
	"Configuration Concepts",
	"Deployment Modes",
	"Kubernetes Resources",
	"Broker Components",
	"Metrics & Monitoring",
	"High Availability",
	"Probes",
	"Validation",
	"Certificate Management",
	"Version Management",
	"Operational Controls",
	"Storage & Persistence",
	"Status & Reporting",
	"Conventions",
}

// RunGlossary executes the glossary generation command
func RunGlossary(args []string) {
	// Parse command-specific flags
	fs := flag.NewFlagSet("glossary", flag.ExitOnError)
	configFile := fs.String("config", "glossary_terms.json", "Terms configuration file")
	outputFile := fs.String("output", "glossary.md", "Output file path")
	quiet := fs.Bool("quiet", false, "Minimal output")
	
	fs.Usage = func() {
		fmt.Println(`Generate glossary from terms configuration

Usage:
  ddt glossary [flags]

Flags:
  --config FILE    Terms config file (default: glossary_terms.json)
  --output FILE    Output file (default: glossary.md)
  --quiet          Minimal output
  --help           Show this help

Examples:
  # Generate glossary with defaults
  ddt glossary

  # Custom config and output
  ddt glossary --config my_terms.json --output custom_glossary.md

  # Quiet mode
  ddt glossary --quiet`)
	}
	
	fs.Parse(args)

	if !*quiet {
		fmt.Println(strings.Repeat("=", 70))
		fmt.Println("Developer Documentation Tools - Glossary Generator")
		fmt.Println(strings.Repeat("=", 70))
		fmt.Println()
	}

	// Get documentation directory
	docDir, err := shared.GetDocDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting documentation directory: %v\n", err)
		os.Exit(1)
	}

	if !*quiet {
		fmt.Printf("Documentation directory: %s\n", docDir)
	}

	// Load terms configuration
	configPath := filepath.Join(docDir, *configFile)
	terms, err := loadTermsConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading terms config: %v\n", err)
		os.Exit(1)
	}

	if !*quiet {
		fmt.Printf("Total terms to process: %d\n\n", len(terms))
	}

	// Verify all documentation files exist
	if !*quiet {
		fmt.Println("Verifying documentation files...")
	}
	
	missingFiles := []string{}
	for _, docFile := range shared.DocFiles {
		fullPath := filepath.Join(docDir, docFile)
		if _, err := os.Stat(fullPath); err == nil {
			if !*quiet {
				fmt.Printf("  ✓ %s\n", docFile)
			}
		} else {
			if !*quiet {
				fmt.Printf("  ✗ %s (NOT FOUND)\n", docFile)
			}
			missingFiles = append(missingFiles, docFile)
		}
	}

	if len(missingFiles) > 0 {
		fmt.Printf("\nERROR: Missing %d documentation file(s)\n", len(missingFiles))
		os.Exit(1)
	}

	if !*quiet {
		fmt.Println("\nGenerating glossary...")
	}

	// Generate glossary content
	glossaryContent := generateGlossaryContent(docDir, terms)

	// Write output file
	outputPath := filepath.Join(docDir, *outputFile)
	err = os.WriteFile(outputPath, []byte(glossaryContent), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing glossary file: %v\n", err)
		os.Exit(1)
	}

	if !*quiet {
		fmt.Println()
		fmt.Println(strings.Repeat("=", 70))
		fmt.Println("✓ Glossary generated successfully!")
		fmt.Printf("  Output: %s\n", outputPath)
		fmt.Printf("  Size: %s bytes\n", shared.FormatNumber(len(glossaryContent)))
		fmt.Println(strings.Repeat("=", 70))
	} else {
		fmt.Printf("Generated: %s (%s bytes)\n", outputPath, shared.FormatNumber(len(glossaryContent)))
	}
}

// loadTermsConfig loads the glossary terms from JSON config file
func loadTermsConfig(filepath string) ([]GlossaryTerm, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config TermsConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return config.Terms, nil
}

// extractSectionHeaders extracts markdown section headers and their anchors from a file
func extractSectionHeaders(filepath string) (map[string]string, error) {
	headers := make(map[string]string)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Pattern for markdown headers (## through ######)
	headerPattern := regexp.MustCompile(`^(#{2,6})\s+(.+)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := headerPattern.FindStringSubmatch(line)
		if matches != nil {
			title := strings.TrimSpace(matches[2])
			anchor := generateAnchor(title)
			headers[anchor] = title
		}
	}

	return headers, scanner.Err()
}

// generateAnchor generates a GitHub-compatible anchor from a section title
func generateAnchor(title string) string {
	// Convert to lowercase
	anchor := strings.ToLower(title)

	// Remove special characters (keep alphanumeric, spaces, and hyphens)
	re := regexp.MustCompile(`[^\w\s-]`)
	anchor = re.ReplaceAllString(anchor, "")

	// Replace spaces and multiple hyphens with single hyphen
	re = regexp.MustCompile(`[-\s]+`)
	anchor = re.ReplaceAllString(anchor, "-")

	// Trim hyphens from ends
	anchor = strings.Trim(anchor, "-")

	return anchor
}

// findTermInFile finds a term in a file and returns relevant section links
func findTermInFile(term GlossaryTerm, filepath string, headers map[string]string) ([]SectionMatch, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	// Check if any search pattern appears in the file
	found := false
	for _, pattern := range term.SearchPatterns {
		// Case-insensitive search
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(pattern))
		if re.MatchString(contentStr) {
			found = true
			break
		}
	}

	if !found {
		return []SectionMatch{}, nil
	}

	// Split content by headers to find relevant sections
	headerPattern := regexp.MustCompile(`^(#{2,6}\s+.+)$`)
	lines := strings.Split(contentStr, "\n")

	var matches []SectionMatch
	currentSection := ""
	currentAnchor := ""
	sectionContent := strings.Builder{}
	seen := make(map[string]bool)

	for _, line := range lines {
		if headerPattern.MatchString(line) {
			// Process previous section
			if currentSection != "" && currentAnchor != "" {
				sectionText := sectionContent.String()
				for _, pattern := range term.SearchPatterns {
					re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(pattern))
					if re.MatchString(sectionText) {
						if !seen[currentAnchor] && headers[currentAnchor] != "" {
							matches = append(matches, SectionMatch{
								Anchor: currentAnchor,
								Title:  currentSection,
							})
							seen[currentAnchor] = true
						}
						break
					}
				}
			}

			// Start new section
			title := strings.TrimSpace(strings.TrimLeft(line, "#"))
			currentSection = title
			currentAnchor = generateAnchor(title)
			sectionContent.Reset()
		} else {
			sectionContent.WriteString(line)
			sectionContent.WriteString("\n")
		}
	}

	// Process last section
	if currentSection != "" && currentAnchor != "" {
		sectionText := sectionContent.String()
		for _, pattern := range term.SearchPatterns {
			re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(pattern))
			if re.MatchString(sectionText) {
				if !seen[currentAnchor] && headers[currentAnchor] != "" {
					matches = append(matches, SectionMatch{
						Anchor: currentAnchor,
						Title:  currentSection,
					})
					seen[currentAnchor] = true
				}
				break
			}
		}
	}

	// Limit to top 3 most relevant sections
	if len(matches) > 3 {
		matches = matches[:3]
	}

	return matches, nil
}

// generateAlphabeticalIndex generates the alphabetical index section
func generateAlphabeticalIndex(terms []GlossaryTerm) string {
	var sb strings.Builder

	// Sort terms alphabetically
	sortedTerms := make([]GlossaryTerm, len(terms))
	copy(sortedTerms, terms)
	sort.Slice(sortedTerms, func(i, j int) bool {
		return strings.ToLower(sortedTerms[i].Name) < strings.ToLower(sortedTerms[j].Name)
	})

	// Group by first letter
	currentLetter := ""
	for i, term := range sortedTerms {
		firstLetter := strings.ToUpper(string([]rune(term.Name)[0]))

		if firstLetter != currentLetter {
			if currentLetter != "" {
				sb.WriteString("\n")
			}
			sb.WriteString("**")
			sb.WriteString(firstLetter)
			sb.WriteString("** ")
			currentLetter = firstLetter
		} else {
			sb.WriteString(" • ")
		}

		anchor := generateAnchor(term.Name)
		sb.WriteString("[")
		sb.WriteString(term.Name)
		sb.WriteString("](#")
		sb.WriteString(anchor)
		sb.WriteString(")")

		// Add newline after last term
		if i == len(sortedTerms)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// generateGlossaryContent generates the complete glossary markdown content
func generateGlossaryContent(docDir string, terms []GlossaryTerm) string {
	var sb strings.Builder

	// Get current date
	today := time.Now().Format("2006-01-02")

	// Build markdown header
	sb.WriteString(`---
title: "Developer Documentation Glossary"
description: "Quick reference index of technical terms used across ActiveMQ Artemis Operator developer documentation"
draft: false
images: []
menu:
  docs:
    parent: "developer"
weight: 299
toc: true
---

> **Auto-Generated Documentation**
> 
> This glossary is automatically generated by ` + "`ddt glossary`" + `.
> Last generated: `)
	sb.WriteString(today)
	sb.WriteString(`
>
> To regenerate: ` + "`cd docs/developer/ddt && go run main.go glossary`" + `

## How to Use This Glossary

This glossary provides quick definitions and navigation links for technical terms used throughout the ActiveMQ Artemis Operator developer documentation. Each term includes:

- **Definition**: A concise explanation of the concept
- **Found in**: Links to relevant sections across the documentation where the term is discussed

Use this as a quick reference when reading the documentation or to understand terminology used in the codebase.

## Documentation Files Indexed

- **[Operator Architecture](operator_architecture.md)**: Technical architecture and implementation details
- **[Operator Conventions](operator_conventions.md)**: Naming patterns, defaults, and automatic behaviors
- **[Operator Defaults Reference](operator_defaults_reference.md)**: Complete default value specifications
- **[Contribution Guide](contribution_guide.md)**: Guide for extending and developing the operator
- **[TDD Test Index](tdd_index.md)**: Comprehensive test catalog (396+ scenarios)

---

## Alphabetical Index

`)

	// Generate alphabetical index
	sb.WriteString(generateAlphabeticalIndex(terms))

	sb.WriteString(`
---

## Terms by Category

`)

	// Extract section headers from all doc files
	fileHeaders := make(map[string]map[string]string)
	for _, docFile := range shared.DocFiles {
		fullPath := filepath.Join(docDir, docFile)
		headers, err := extractSectionHeaders(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to extract headers from %s: %v\n", docFile, err)
			continue
		}
		fileHeaders[docFile] = headers
	}

	// Group terms by category
	categories := make(map[string][]GlossaryTerm)
	for _, term := range terms {
		categories[term.Category] = append(categories[term.Category], term)
	}

	// Generate content for each category in order
	for _, category := range categoryOrder {
		termsInCategory, exists := categories[category]
		if !exists {
			continue
		}

		sb.WriteString("### ")
		sb.WriteString(category)
		sb.WriteString("\n\n")

		// Sort terms alphabetically within category
		sort.Slice(termsInCategory, func(i, j int) bool {
			return termsInCategory[i].Name < termsInCategory[j].Name
		})

		for _, term := range termsInCategory {
			sb.WriteString("#### ")
			sb.WriteString(term.Name)
			sb.WriteString("\n\n")
			sb.WriteString(term.Definition)
			sb.WriteString("\n\n")

			// Find occurrences in documentation
			var foundIn []string
			for _, docFile := range shared.DocFiles {
				fullPath := filepath.Join(docDir, docFile)
				sections, err := findTermInFile(term, fullPath, fileHeaders[docFile])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: error finding term in %s: %v\n", docFile, err)
					continue
				}

				for _, section := range sections {
					// Create clean doc name for display
					docName := strings.TrimSuffix(docFile, ".md")
					docName = strings.ReplaceAll(docName, "_", " ")
					docName = strings.Title(docName)

					link := fmt.Sprintf("- [%s - %s](%s#%s)",
						docName, section.Title, docFile, section.Anchor)
					foundIn = append(foundIn, link)
				}
			}

			if len(foundIn) > 0 {
				sb.WriteString("**Found in:**\n")
				for _, link := range foundIn {
					sb.WriteString(link)
					sb.WriteString("\n")
				}
				sb.WriteString("\n")
			} else {
				sb.WriteString("*Not explicitly referenced in current documentation.*\n\n")
			}

			sb.WriteString("---\n\n")
		}
	}

	// Add footer with stats
	totalTerms := len(terms)
	totalCategories := len(categories)

	sb.WriteString(fmt.Sprintf(`
## Glossary Statistics

- **Total Terms**: %d
- **Categories**: %d
- **Documentation Files**: %d

---

## Adding New Terms

To add terms to this glossary, edit `+"`glossary_terms.json`"+` and add entries to the `+"`terms`"+` array:

`+"```json"+`
{
  "name": "Your Term",
  "definition": "Clear, concise definition...",
  "searchPatterns": ["term", "alternate-term", "alias"],
  "category": "Appropriate Category"
}
`+"```"+`

Then regenerate with: `+"`cd ddt && go run main.go glossary`"+`

## Related Documentation

- **[Operator Architecture](operator_architecture.md)**: Detailed architecture documentation
- **[Contribution Guide](contribution_guide.md)**: How to contribute to the operator
- **[TDD Test Index](tdd_index.md)**: Complete test catalog
`, totalTerms, totalCategories, len(shared.DocFiles)))

	return sb.String()
}
