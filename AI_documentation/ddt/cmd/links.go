package cmd

import (
	"bytes"

	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/arkmq-org/ddt/shared"

	"github.com/pterm/pterm"
)

const (
	DefaultConfidenceThreshold = 0.95
	DefaultWarningThreshold    = 0.70
	ConfigFileName             = "link_updater_config.json"
)

// DocLink represents a GitHub permalink in documentation
type DocLink struct {
	FilePath   string // Documentation file path
	LineNumber int    // Line number in doc file
	URL        string // Original URL
	CommitSHA  string // Extracted commit SHA
	RepoOwner  string // Repository owner
	RepoName   string // Repository name
	TargetPath string // Source code file path
	StartLine  int    // Line number or range start
	EndLine    int    // Range end (0 if single line)
	LinkText   string // Markdown link text
	FullMatch  string // Full markdown link syntax
}

// UpdateResult represents the result of processing a link
type UpdateResult struct {
	Link       DocLink
	Status     string  // "updated", "warning", "error", "skipped"
	NewSHA     string  // Target commit SHA
	NewStart   int     // New line number
	NewEnd     int     // New end line (if range)
	NewURL     string  // Updated URL
	Confidence float64 // Match confidence (0.0-1.0)
	Message    string  // Human-readable explanation
}

// MatchResult represents a code matching result
type MatchResult struct {
	Found      bool
	LineNumber int
	EndLine    int
	Confidence float64
	Reason     string
}

// Config holds the tool configuration
type Config struct {
	TargetCommit       string
	DocumentationPaths []string
	UpdateThreshold    float64
	WarningThreshold   float64
	CreateBackups      bool
}

var (
	// GitHub permalink pattern
	githubLinkPattern = regexp.MustCompile(
		`\[([^\]]+)\]\((https://github\.com/([^/]+)/([^/]+)/blob/([0-9a-f]{7,40})/([^#\)]+)(?:#L(\d+)(?:-L(\d+))?)?)\)`,
	)

	// Global statistics
	stats = struct {
		Total    int
		Updated  int
		Warnings int
		Errors   int
		Skipped  int
	}{}

	// Command flags (initialized in RunLinks)
	updateFlag      *bool
	interactiveFlag *bool
	reportFlag      *bool
	confidenceFlag  *float64
	outputFlag      *string
	quietFlag       *bool
	noStyleFlag     *bool
	verboseFlag     *bool
	configFlag      *string
)

func RunLinks(args []string) {
	// Parse command-specific flags
	fs := flag.NewFlagSet("links", flag.ExitOnError)

	updateFlag = fs.Bool("update", false, "Apply updates to documentation files")
	interactiveFlag = fs.Bool("interactive", false, "Prompt for each low-confidence update")
	reportFlag = fs.Bool("report", false, "Generate markdown report")
	confidenceFlag = fs.Float64("confidence", DefaultConfidenceThreshold, "Confidence threshold for auto-update")
	outputFlag = fs.String("output", "", "Report output file")
	quietFlag = fs.Bool("quiet", false, "Minimal output")
	noStyleFlag = fs.Bool("no-style", false, "Disable pterm styling")
	verboseFlag = fs.Bool("verbose", false, "Detailed logging")
	configFlag = fs.String("config", ConfigFileName, "Config file path")

	fs.Parse(args)

	// Get target commit from remaining args
	targetCommit := "HEAD"
	if fs.NArg() > 0 {
		targetCommit = fs.Arg(0)
	}

	// Configure pterm
	if *quietFlag || *noStyleFlag {
		pterm.DisableStyling()
		pterm.DisableColor()
	}

	// Display header
	if !*quietFlag {
		pterm.DefaultHeader.WithFullWidth().Println("Documentation Link Updater")
		pterm.Println()
	}

	// Find documentation files
	docDir, err := shared.GetDocDir()
	if err != nil {
		pterm.Error.Printf("Failed to get documentation directory: %v\n", err)
		os.Exit(1)
	}

	docFiles, err := shared.FindDocFiles(docDir)
	if err != nil {
		pterm.Error.Printf("Failed to find documentation files: %v\n", err)
		os.Exit(1)
	}

	if *verboseFlag {
		pterm.Info.Printf("Found %d documentation files\n", len(docFiles))
	}

	// Phase 1: Discover links
	if !*quietFlag {
		pterm.DefaultSection.Println("Phase 1: Discovering Links")
	}

	links, sourceCommit, err := discoverLinks(docFiles)
	if err != nil {
		pterm.Error.Printf("Failed to discover links: %v\n", err)
		os.Exit(1)
	}

	if len(links) == 0 {
		pterm.Warning.Println("No GitHub permalinks found in documentation")
		return
	}

	stats.Total = len(links)

	// Display detection results
	if !*quietFlag {
		pterm.Success.Printf("Detected source commit: %s (from %d links)\n", sourceCommit[:8], len(links))
	}

	// Resolve target commit to full SHA
	targetSHA, err := resolveCommit(targetCommit)
	if err != nil {
		pterm.Error.Printf("Failed to resolve target commit: %v\n", err)
		os.Exit(1)
	}

	if !*quietFlag {
		pterm.Info.Printf("Target commit: %s (%s)\n", targetCommit, targetSHA[:8])
		pterm.Println()
	}

	// Check if source and target are the same
	if sourceCommit == targetSHA {
		pterm.Info.Println("Source and target commits are the same - nothing to update")
		return
	}

	// Phase 2: Analyze links
	if !*quietFlag {
		pterm.DefaultSection.Println("Phase 2: Analyzing Links")
	}

	results := analyzeLinks(links, sourceCommit, targetSHA)

	// Categorize results
	for _, result := range results {
		switch result.Status {
		case "updated":
			stats.Updated++
		case "warning":
			stats.Warnings++
		case "error":
			stats.Errors++
		case "skipped":
			stats.Skipped++
		}
	}

	// Phase 3: Display results
	displaySummary(results)

	// Phase 4: Apply updates if requested
	if *updateFlag {
		if !*quietFlag {
			pterm.Println()
			pterm.DefaultSection.Println("Phase 3: Applying Updates")
		}

		// Confirm if not in quiet mode
		if !*quietFlag && !*interactiveFlag {
			confirm, _ := pterm.DefaultInteractiveConfirm.
				WithDefaultText(fmt.Sprintf("Apply updates to %d file(s)?", len(docFiles))).
				WithDefaultValue(false).
				Show()

			if !confirm {
				pterm.Info.Println("Update cancelled")
				return
			}
		}

		err := applyUpdates(results, *interactiveFlag)
		if err != nil {
			pterm.Error.Printf("Failed to apply updates: %v\n", err)
			os.Exit(1)
		}

		if !*quietFlag {
			pterm.Success.Println("Updates applied successfully")
		}
	}

	// Generate report if requested
	if *reportFlag {
		outputPath := *outputFlag
		if outputPath == "" {
			outputPath = "link_update_report.md"
		}

		err := generateReport(results, sourceCommit, targetSHA, outputPath)
		if err != nil {
			pterm.Error.Printf("Failed to generate report: %v\n", err)
			os.Exit(1)
		}

		if !*quietFlag {
			pterm.Success.Printf("Report generated: %s\n", outputPath)
		}
	}

	// Exit with non-zero if there are warnings or errors
	if stats.Warnings > 0 || stats.Errors > 0 {
		os.Exit(1)
	}
}

// discoverLinks scans documentation files and extracts GitHub permalinks
func discoverLinks(docFiles []string) ([]DocLink, string, error) {
	var allLinks []DocLink
	commitCounts := make(map[string]int)

	var spinner *pterm.SpinnerPrinter
	if !*quietFlag {
		spinner, _ = pterm.DefaultSpinner.Start("Scanning documentation files...")
	}

	for _, docFile := range docFiles {
		content, err := os.ReadFile(docFile)
		if err != nil {
			if spinner != nil {
				spinner.Fail(fmt.Sprintf("Failed to read %s", docFile))
			}
			return nil, "", err
		}

		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			matches := githubLinkPattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				link := DocLink{
					FilePath:   docFile,
					LineNumber: lineNum + 1,
					LinkText:   match[1],
					URL:        match[2],
					RepoOwner:  match[3],
					RepoName:   match[4],
					CommitSHA:  match[5],
					TargetPath: match[6],
					FullMatch:  match[0],
				}

				// Parse line numbers
				if match[7] != "" {
					fmt.Sscanf(match[7], "%d", &link.StartLine)
					if match[8] != "" {
						fmt.Sscanf(match[8], "%d", &link.EndLine)
					}
				}

				allLinks = append(allLinks, link)
				commitCounts[link.CommitSHA]++
			}
		}

		if *verboseFlag && spinner != nil {
			spinner.UpdateText(fmt.Sprintf("Scanned %s: found %d links", docFile, len(allLinks)))
		}
	}

	if spinner != nil {
		spinner.Success(fmt.Sprintf("Found %d links in %d files", len(allLinks), len(docFiles)))
	}

	// Determine most common commit (source commit)
	var sourceCommit string
	maxCount := 0
	for commit, count := range commitCounts {
		if count > maxCount {
			maxCount = count
			sourceCommit = commit
		}
	}

	// Warn if multiple commits detected
	if len(commitCounts) > 1 && !*quietFlag {
		pterm.Warning.Printf("Multiple source commits detected: %d different commits\n", len(commitCounts))
		if *verboseFlag {
			for commit, count := range commitCounts {
				pterm.Printf("  %s: %d links\n", commit[:8], count)
			}
		}
	}

	return allLinks, sourceCommit, nil
}

// resolveCommit resolves a commit reference to full SHA
func resolveCommit(ref string) (string, error) {
	cmd := exec.Command("git", "rev-parse", ref)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to resolve %s: %w", ref, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// analyzeLinks analyzes each link and determines if it can be updated
func analyzeLinks(links []DocLink, sourceCommit, targetCommit string) []UpdateResult {
	results := make([]UpdateResult, 0, len(links))

	var pb *pterm.ProgressbarPrinter
	if !*quietFlag {
		pb, _ = pterm.DefaultProgressbar.
			WithTotal(len(links)).
			WithTitle("Processing links").
			Start()
	}

	for _, link := range links {
		result := processLink(link, sourceCommit, targetCommit)
		results = append(results, result)

		if pb != nil {
			pb.Increment()
		}
	}

	return results
}

// processLink processes a single link and returns the result
func processLink(link DocLink, sourceCommit, targetCommit string) UpdateResult {
	result := UpdateResult{
		Link:   link,
		NewSHA: targetCommit,
	}

	// Check if file exists in target commit
	if !fileExistsInCommit(targetCommit, link.TargetPath) {
		result.Status = "error"
		result.Message = "File not found in target commit"
		return result
	}

	// For links without line numbers, just update the commit SHA
	if link.StartLine == 0 {
		result.Status = "updated"
		result.Confidence = 1.0
		result.Message = "File exists (no line number)"
		result.NewURL = strings.Replace(link.URL, link.CommitSHA, targetCommit, 1)
		return result
	}

	// Get file content at both commits
	sourceContent, err := getFileContent(sourceCommit, link.TargetPath)
	if err != nil {
		result.Status = "error"
		result.Message = fmt.Sprintf("Failed to get source content: %v", err)
		return result
	}

	targetContent, err := getFileContent(targetCommit, link.TargetPath)
	if err != nil {
		result.Status = "error"
		result.Message = fmt.Sprintf("Failed to get target content: %v", err)
		return result
	}

	// Extract target lines from source
	sourceLines := strings.Split(sourceContent, "\n")
	if link.StartLine > len(sourceLines) {
		result.Status = "error"
		result.Message = "Line number out of range in source"
		return result
	}

	endLine := link.StartLine
	if link.EndLine > 0 {
		endLine = link.EndLine
	}

	targetLines := extractLines(sourceLines, link.StartLine, endLine)

	// Find best match in target
	match := findBestMatch(targetLines, targetContent)

	if !match.Found {
		result.Status = "error"
		result.Confidence = 0.0
		result.Message = match.Reason
		return result
	}

	result.NewStart = match.LineNumber
	result.NewEnd = match.EndLine
	result.Confidence = match.Confidence

	// Build new URL
	newURL := strings.Replace(link.URL, link.CommitSHA, targetCommit, 1)
	if link.StartLine > 0 {
		// Replace line numbers
		oldFragment := fmt.Sprintf("#L%d", link.StartLine)
		newFragment := fmt.Sprintf("#L%d", result.NewStart)
		if link.EndLine > 0 {
			oldFragment = fmt.Sprintf("#L%d-L%d", link.StartLine, link.EndLine)
			newFragment = fmt.Sprintf("#L%d-L%d", result.NewStart, result.NewEnd)
		}
		newURL = strings.Replace(newURL, oldFragment, newFragment, 1)
	}
	result.NewURL = newURL

	// Classify based on confidence
	if result.Confidence >= *confidenceFlag {
		result.Status = "updated"
		result.Message = "High confidence match"
	} else if result.Confidence >= DefaultWarningThreshold {
		result.Status = "warning"
		result.Message = fmt.Sprintf("Low confidence match (%.2f)", result.Confidence)
	} else {
		result.Status = "error"
		result.Message = fmt.Sprintf("Very low confidence (%.2f)", result.Confidence)
	}

	return result
}

// fileExistsInCommit checks if a file exists in a specific commit
func fileExistsInCommit(commit, path string) bool {
	cmd := exec.Command("git", "ls-tree", commit, path)
	return cmd.Run() == nil
}

// getFileContent retrieves file content at a specific commit
func getFileContent(commit, path string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", commit, path))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// extractLines extracts a range of lines from content
func extractLines(lines []string, start, end int) []string {
	if start < 1 || start > len(lines) {
		return []string{}
	}
	if end == 0 || end < start {
		end = start
	}
	if end > len(lines) {
		end = len(lines)
	}
	return lines[start-1 : end]
}

// findBestMatch finds the best matching location for target lines in content
func findBestMatch(targetLines []string, content string) MatchResult {
	lines := strings.Split(content, "\n")

	// Strategy 1: Exact match
	if match := findExactMatch(targetLines, lines); match.Found {
		return match
	}

	// Strategy 2: Normalized match (ignore whitespace differences)
	if match := findNormalizedMatch(targetLines, lines); match.Found {
		return match
	}

	// Strategy 3: Fuzzy match with context
	if match := findFuzzyMatch(targetLines, lines); match.Found {
		return match
	}

	return MatchResult{
		Found:  false,
		Reason: "Content not found or significantly changed",
	}
}

// findExactMatch finds exact content match
func findExactMatch(targetLines, contentLines []string) MatchResult {
	targetStr := strings.Join(targetLines, "\n")

	for i := 0; i <= len(contentLines)-len(targetLines); i++ {
		candidateStr := strings.Join(contentLines[i:i+len(targetLines)], "\n")
		if candidateStr == targetStr {
			return MatchResult{
				Found:      true,
				LineNumber: i + 1,
				EndLine:    i + len(targetLines),
				Confidence: 1.0,
				Reason:     "Exact match",
			}
		}
	}

	return MatchResult{Found: false}
}

// findNormalizedMatch finds match ignoring whitespace
func findNormalizedMatch(targetLines, contentLines []string) MatchResult {
	normalize := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	targetNorm := make([]string, len(targetLines))
	for i, line := range targetLines {
		targetNorm[i] = normalize(line)
	}
	targetStr := strings.Join(targetNorm, "\n")

	for i := 0; i <= len(contentLines)-len(targetLines); i++ {
		candidateNorm := make([]string, len(targetLines))
		for j := 0; j < len(targetLines); j++ {
			candidateNorm[j] = normalize(contentLines[i+j])
		}
		candidateStr := strings.Join(candidateNorm, "\n")

		if candidateStr == targetStr {
			return MatchResult{
				Found:      true,
				LineNumber: i + 1,
				EndLine:    i + len(targetLines),
				Confidence: 0.95,
				Reason:     "Normalized match",
			}
		}
	}

	return MatchResult{Found: false}
}

// findFuzzyMatch finds match with some tolerance
func findFuzzyMatch(targetLines, contentLines []string) MatchResult {
	if len(targetLines) == 0 {
		return MatchResult{Found: false}
	}

	// Look for the first line with high similarity
	firstLine := strings.TrimSpace(targetLines[0])
	bestMatch := MatchResult{Found: false}
	bestSimilarity := 0.0

	for i, line := range contentLines {
		sim := stringSimilarity(firstLine, strings.TrimSpace(line))
		if sim > 0.8 && sim > bestSimilarity {
			// Check if subsequent lines also match
			totalSim := sim
			matches := 1

			for j := 1; j < len(targetLines) && i+j < len(contentLines); j++ {
				lineSim := stringSimilarity(
					strings.TrimSpace(targetLines[j]),
					strings.TrimSpace(contentLines[i+j]),
				)
				totalSim += lineSim
				matches++
			}

			avgSim := totalSim / float64(matches)
			if avgSim > 0.75 && avgSim > bestSimilarity {
				bestSimilarity = avgSim
				bestMatch = MatchResult{
					Found:      true,
					LineNumber: i + 1,
					EndLine:    i + len(targetLines),
					Confidence: avgSim * 0.9, // Reduce confidence for fuzzy match
					Reason:     "Fuzzy match",
				}
			}
		}
	}

	return bestMatch
}

// stringSimilarity calculates similarity between two strings (0.0-1.0)
func stringSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	// Simple similarity: ratio of matching characters
	s1Runes := []rune(s1)
	s2Runes := []rune(s2)

	minLen := len(s1Runes)
	maxLen := len(s2Runes)
	if minLen > maxLen {
		minLen, maxLen = maxLen, minLen
	}

	if maxLen == 0 {
		return 0.0
	}

	matches := 0
	for i := 0; i < minLen; i++ {
		if s1Runes[i] == s2Runes[i] {
			matches++
		}
	}

	return float64(matches) / float64(maxLen)
}

// displaySummary displays a summary of the analysis results
func displaySummary(results []UpdateResult) {
	if *quietFlag {
		fmt.Printf("Source: %s\n", results[0].Link.CommitSHA[:8])
		fmt.Printf("Target: %s\n", results[0].NewSHA[:8])
		fmt.Printf("Total: %d\n", stats.Total)
		fmt.Printf("Updated: %d\n", stats.Updated)
		fmt.Printf("Warnings: %d\n", stats.Warnings)
		fmt.Printf("Errors: %d\n", stats.Errors)
		return
	}

	pterm.Println()
	pterm.DefaultSection.Println("Summary")

	// Create summary table
	tableData := pterm.TableData{
		{"Status", "Count", "Percentage"},
		{pterm.Green("Can Update"), fmt.Sprintf("%d", stats.Updated), fmt.Sprintf("%.1f%%", float64(stats.Updated)/float64(stats.Total)*100)},
		{pterm.Yellow("Warnings"), fmt.Sprintf("%d", stats.Warnings), fmt.Sprintf("%.1f%%", float64(stats.Warnings)/float64(stats.Total)*100)},
		{pterm.Red("Errors"), fmt.Sprintf("%d", stats.Errors), fmt.Sprintf("%.1f%%", float64(stats.Errors)/float64(stats.Total)*100)},
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	// Group by file
	if *verboseFlag {
		pterm.Println()
		pterm.Info.Println("Per-file breakdown:")

		fileStats := make(map[string]map[string]int)
		for _, result := range results {
			if fileStats[result.Link.FilePath] == nil {
				fileStats[result.Link.FilePath] = make(map[string]int)
			}
			fileStats[result.Link.FilePath][result.Status]++
		}

		for file, counts := range fileStats {
			pterm.Printf("  %s:\n", file)
			pterm.Printf("    %s %d   %s %d   %s %d\n",
				pterm.Green("✓"), counts["updated"],
				pterm.Yellow("⚠"), counts["warning"],
				pterm.Red("✗"), counts["error"])
		}
	}
}

// applyUpdates applies the updates to documentation files
func applyUpdates(results []UpdateResult, interactive bool) error {
	// Group by file
	fileUpdates := make(map[string][]UpdateResult)
	for _, result := range results {
		if result.Status == "updated" || (interactive && result.Status == "warning") {
			fileUpdates[result.Link.FilePath] = append(fileUpdates[result.Link.FilePath], result)
		}
	}

	for filePath, updates := range fileUpdates {
		// Read file
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filePath, err)
		}

		modified := string(content)
		appliedCount := 0

		// Apply updates (in reverse order to maintain line numbers)
		for i := len(updates) - 1; i >= 0; i-- {
			result := updates[i]

			// Interactive confirmation for warnings
			if interactive && result.Status == "warning" {
				pterm.Println()
				pterm.Warning.Printf("Low confidence update (%.2f):\n", result.Confidence)
				pterm.Printf("  File: %s:%d\n", result.Link.FilePath, result.Link.LineNumber)
				pterm.Printf("  Old: %s\n", result.Link.URL)
				pterm.Printf("  New: %s\n", result.NewURL)

				confirm, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText("Apply this update?").
					WithDefaultValue(false).
					Show()

				if !confirm {
					continue
				}
			}

			// Replace the old URL with new URL
			modified = strings.Replace(modified, result.Link.FullMatch,
				strings.Replace(result.Link.FullMatch, result.Link.URL, result.NewURL, 1), 1)
			appliedCount++
		}

		if appliedCount > 0 {
			// Write back
			err = os.WriteFile(filePath, []byte(modified), 0644)
			if err != nil {
				return fmt.Errorf("failed to write %s: %w", filePath, err)
			}

			if !*quietFlag {
				pterm.Success.Printf("Updated %s (%d changes)\n", filePath, appliedCount)
			}
		}
	}

	return nil
}

// generateReport generates a markdown report
func generateReport(results []UpdateResult, sourceCommit, targetCommit, outputPath string) error {
	var buf bytes.Buffer

	buf.WriteString("# Documentation Link Update Report\n\n")
	buf.WriteString(fmt.Sprintf("Generated: %s\n\n", pterm.Sprintf("2025-01-10")))
	buf.WriteString(fmt.Sprintf("- **Source Commit**: `%s`\n", sourceCommit))
	buf.WriteString(fmt.Sprintf("- **Target Commit**: `%s`\n\n", targetCommit))

	buf.WriteString("## Summary\n\n")
	buf.WriteString(fmt.Sprintf("- Total Links: %d\n", stats.Total))
	buf.WriteString(fmt.Sprintf("- ✓ Can Update: %d (%.1f%%)\n", stats.Updated, float64(stats.Updated)/float64(stats.Total)*100))
	buf.WriteString(fmt.Sprintf("- ⚠ Warnings: %d (%.1f%%)\n", stats.Warnings, float64(stats.Warnings)/float64(stats.Total)*100))
	buf.WriteString(fmt.Sprintf("- ✗ Errors: %d (%.1f%%)\n\n", stats.Errors, float64(stats.Errors)/float64(stats.Total)*100))

	// Group results by status and file
	for _, status := range []string{"updated", "warning", "error"} {
		fileResults := make(map[string][]UpdateResult)
		for _, result := range results {
			if result.Status == status {
				fileResults[result.Link.FilePath] = append(fileResults[result.Link.FilePath], result)
			}
		}

		if len(fileResults) == 0 {
			continue
		}

		switch status {
		case "updated":
			buf.WriteString("## ✓ Can Update (High Confidence)\n\n")
		case "warning":
			buf.WriteString("## ⚠ Warnings (Manual Review Recommended)\n\n")
		case "error":
			buf.WriteString("## ✗ Errors (Manual Fix Required)\n\n")
		}

		for file, results := range fileResults {
			buf.WriteString(fmt.Sprintf("### %s\n\n", file))
			for _, result := range results {
				buf.WriteString(fmt.Sprintf("**Line %d**: ", result.Link.LineNumber))
				if result.Status == "updated" {
					buf.WriteString(fmt.Sprintf("`%s#L%d` → `#L%d`\n",
						result.Link.TargetPath, result.Link.StartLine, result.NewStart))
					buf.WriteString(fmt.Sprintf("  - Confidence: %.2f\n", result.Confidence))
				} else {
					buf.WriteString(fmt.Sprintf("`%s#L%d`\n", result.Link.TargetPath, result.Link.StartLine))
					buf.WriteString(fmt.Sprintf("  - %s\n", result.Message))
					if result.Confidence > 0 {
						buf.WriteString(fmt.Sprintf("  - Confidence: %.2f\n", result.Confidence))
					}
				}
				buf.WriteString("\n")
			}
		}
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}
