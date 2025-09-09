# Developer Documentation Tools (ddt)

A unified CLI tool for maintaining ActiveMQ Artemis Operator developer documentation.

## Features

**ddt** combines two essential documentation maintenance tools into a single, easy-to-use CLI:

1. **Glossary Generation** (`ddt glossary`) - Automatically generates a comprehensive glossary of technical terms from a curated list, with cross-references to relevant documentation sections.

2. **Link Updater** (`ddt links`) - Intelligently updates GitHub permalinks in documentation as the codebase evolves, using git diff analysis and fuzzy matching.

## Installation

### Build from Source

```bash
cd AI_documentation/ddt
go build -o ddt
```

The binary will be created in the current directory.

### Requirements

- Go 1.23 or later
- Git (for the `links` command)

## Commands

### `ddt glossary` - Generate Glossary

Generates a markdown glossary from a curated list of terms, scanning all developer documentation files to find relevant sections where each term is discussed.

**Usage:**

```bash
# Generate glossary with defaults
ddt glossary

# Custom config and output paths
ddt glossary --config ../glossary_terms.json --output ../glossary.md

# Quiet mode (minimal output)
ddt glossary --quiet
```

**Flags:**

- `--config FILE` - Terms configuration file (default: `glossary_terms.json`)
- `--output FILE` - Output file path (default: `glossary.md`)
- `--quiet` - Minimal output
- `--help` - Show command help

**Configuration:**

Terms are defined in `glossary_terms.json` (located in the parent directory):

```json
{
  "terms": [
    {
      "name": "Term Name",
      "definition": "Clear, concise definition",
      "searchPatterns": ["term", "alternate-term"],
      "category": "Category Name"
    }
  ]
}
```

**Output:**

- Alphabetical index of all terms
- Terms organized by category
- Cross-references to documentation sections
- Statistics (total terms, categories, files)

### `ddt links` - Update Documentation Links

Automatically updates GitHub permalinks in documentation files as the codebase evolves. Uses intelligent matching strategies (exact, normalized, fuzzy) to track code changes across commits.

**Usage:**

```bash
# Check what needs updating (dry-run, shows summary)
ddt links

# Check against specific commit
ddt links v1.2.3

# Apply updates to current HEAD
ddt links --update

# Interactive mode (prompts for low-confidence updates)
ddt links --update --interactive

# Generate detailed markdown report
ddt links --report --output ../link_report.md

# Quiet mode for CI/CD pipelines
ddt links --quiet
```

**Flags:**

- `--update` - Apply updates to documentation files (default: dry-run)
- `--interactive` - Prompt for each low-confidence update
- `--report` - Generate markdown report
- `--confidence FLOAT` - Confidence threshold for auto-update (default: 0.95)
- `--output FILE` - Report output file path
- `--quiet` - Minimal output (summary only)
- `--no-style` - Disable pterm styling (for CI/CD)
- `--verbose` - Detailed logging
- `--config FILE` - Config file path (default: `link_updater_config.json`)
- `--help` - Show command help

**How It Works:**

1. **Auto-detects source commit** from existing links in documentation
2. **Analyzes git diff** between source and target commits
3. **Applies intelligent matching**:
   - Exact match (1.0 confidence)
   - Normalized match (0.95 confidence) - ignores whitespace/comments
   - Fuzzy match (0.75-0.90 confidence) - similar code structure
4. **Reports results** with confidence scores
5. **Updates documentation** (if `--update` flag is used)

**Matching Strategies:**

- **Exact Match** (1.0): Code is identical
- **Normalized Match** (0.95): Same code, different formatting
- **Fuzzy Match** (0.75-0.90): Similar structure, some changes

**Configuration:**

Settings are stored in `link_updater_config.json` (parent directory):

```json
{
  "updateThreshold": 0.95,
  "warningThreshold": 0.70
}
```

**CI/CD Integration:**

```bash
# In your CI pipeline, check link health
ddt links --quiet --no-style

# Exit code 0: all links OK or high-confidence updates available
# Exit code 1: errors or low-confidence matches need manual review
```

## Common Workflows

### After Code Changes

When you modify source code, update documentation links:

```bash
cd AI_documentation/ddt

# 1. Check what needs updating
./ddt links

# 2. Review the summary

# 3. Apply updates
./ddt links --update

# 4. Review any errors/warnings manually
./ddt links --report --output report.md

# 5. Commit updated documentation
cd ..
git add *.md
git commit -m "docs: update code links to latest commit"
```

### Adding Glossary Terms

1. Edit `glossary_terms.json` (parent directory):

```json
{
  "name": "Your Term",
  "definition": "Clear definition of the concept",
  "searchPatterns": ["term", "alternate-name"],
  "category": "Appropriate Category"
}
```

2. Regenerate glossary:

```bash
cd AI_documentation/ddt
./ddt glossary
```

3. Commit changes:

```bash
cd ..
git add glossary_terms.json glossary.md
git commit -m "docs: add glossary term for XYZ"
```

### Periodic Maintenance

Check documentation health periodically:

```bash
# Generate comprehensive report
./ddt links --report --output health_report.md

# Review warnings and errors
cat health_report.md

# Update high-confidence links
./ddt links --update

# Manually fix remaining issues
```

## Documentation Files Processed

The tool processes these developer documentation files:

- `operator_architecture.md` - Architecture and implementation details
- `operator_conventions.md` - Naming patterns and automatic behaviors
- `operator_defaults_reference.md` - Default value specifications
- `contribution_guide.md` - Development and contribution guide
- `tdd_index.md` - Comprehensive test catalog (396+ scenarios)

## Output Examples

### Glossary Generation

```
$ ./ddt glossary
======================================================================
Developer Documentation Tools - Glossary Generator
======================================================================

Documentation directory: /path/to/AI_documentation
Total terms to process: 112

Verifying documentation files...
  ✓ operator_architecture.md
  ✓ operator_conventions.md
  ✓ operator_defaults_reference.md
  ✓ contribution_guide.md
  ✓ tdd_index.md

Generating glossary...

======================================================================
✓ Glossary generated successfully!
  Output: /path/to/AI_documentation/glossary.md
  Size: 121,433 bytes
======================================================================
```

### Link Update (Dry-Run)

```
$ ./ddt links

╔═══════════════════════════════════════════╗
║      Documentation Link Updater           ║
╚═══════════════════════════════════════════╝

Found 5 documentation files
Scanning documentation files... ✓

Found 1,352 links across 5 files

Source Commit: 4acadb95
Target Commit: 5e9449be

Processing links... ████████████████████ 100%

Results:
  Total:    1,352
  Updated:  1,350
  Warnings: 0
  Errors:   2

Run with --update to apply changes
Run with --report to generate detailed report
```

### Link Update (With --update)

```
$ ./ddt links --update

... (same as above) ...

Applying updates... ████████████████████ 100%

✓ Updated 1,350 links
  operator_architecture.md: 481 updated
  tdd_index.md: 369 updated
  contribution_guide.md: 29 updated
  operator_defaults_reference.md: 34 updated
  operator_conventions.md: 44 updated
  glossary.md: 393 updated

⚠ Manual review needed:
  - operator_architecture.md:1234 (confidence: 0.68)
  - tdd_index.md:5678 (line not found)

Use --report for detailed analysis
```

## Troubleshooting

### Glossary Issues

**Problem:** "Error loading terms config"

```bash
# Ensure glossary_terms.json exists in parent directory
ls ../glossary_terms.json

# If running from different directory, use --config
./ddt glossary --config /path/to/glossary_terms.json
```

**Problem:** "Missing documentation file(s)"

```bash
# Verify you're in the correct directory
pwd  # Should be AI_documentation/ddt/

# Check that documentation files exist in parent directory
ls ../operator_architecture.md
```

### Link Updater Issues

**Problem:** "Could not detect source commit"

- Ensure documentation files contain existing GitHub permalinks
- Check that links follow the format: `https://github.com/owner/repo/blob/SHA/file#L123`

**Problem:** "Many low-confidence matches"

- Large refactorings may cause fuzzy matching
- Use `--interactive` to review each change
- Or use `--report` to analyze before updating

**Problem:** Git errors

```bash
# Ensure you're in a git repository
git status

# Verify target commit exists
git rev-parse HEAD
```

## Development

### Project Structure

```
ddt/
├── main.go                   # Entry point, command dispatcher
├── cmd/
│   ├── glossary.go          # Glossary generation logic
│   └── links.go             # Link update logic
├── shared/
│   └── common.go            # Shared utilities
├── go.mod                    # Module definition
├── go.sum                    # Dependencies lock
├── README.md                 # This file
├── LICENSE                   # Apache 2.0 license
└── .gitignore               # Git ignore rules
```

### Dependencies

- **pterm** (v0.12.81) - Rich terminal output and styling
- Standard library only otherwise

### Building

```bash
# Development build
go build -o ddt

# Optimized production build
go build -ldflags="-s -w" -o ddt

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o ddt-linux

# Run tests (when added)
go test ./...
```

### Contributing

This tool is part of the ActiveMQ Artemis Operator project. Contributions are welcome!

1. Test your changes thoroughly
2. Update this README if adding features
3. Ensure both commands work correctly
4. Follow Go best practices

## License

Apache License 2.0 - See [LICENSE](LICENSE) file for details.

## Related Documentation

- **Parent Directory:** [`AI_documentation/`](../) - All developer documentation
- **Glossary Terms:** [`glossary_terms.json`](../glossary_terms.json) - Term definitions
- **Link Config:** [`link_updater_config.json`](../link_updater_config.json) - Link updater settings
- **Generated Glossary:** [`glossary.md`](../glossary.md) - Auto-generated glossary

## Support

For issues, questions, or contributions:

- **GitHub Issues:** https://github.com/artemiscloud/activemq-artemis-operator/issues
- **Documentation:** https://artemiscloud.io/
- **Community:** https://artemiscloud.io/community/

---

**Note:** This tool was created specifically for the ActiveMQ Artemis Operator developer documentation but could be adapted for other projects with similar needs.

