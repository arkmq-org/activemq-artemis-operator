# DDT Quick Start Guide

## Installation

```bash
cd docs/developer/ddt
go build -o ddt
```

## Common Commands

### Glossary

```bash
# Generate glossary
./ddt glossary

# Quiet mode
./ddt glossary --quiet
```

### Links

```bash
# Check links (dry-run)
./ddt links

# Update links
./ddt links --update

# Generate report
./ddt links --report --output report.md

# CI mode
./ddt links --quiet --no-style
```

## Daily Workflow

```bash
cd docs/developer/ddt

# 1. Check documentation health
./ddt links

# 2. Update if needed
./ddt links --update

# 3. Regenerate glossary
./ddt glossary

# 4. Commit changes
cd ..
git add *.md
git commit -m "docs: update documentation"
```

## Adding Terms

1. Edit `../glossary_terms.json`
2. Run `./ddt glossary`
3. Review `../glossary.md`
4. Commit both files

## Help

```bash
./ddt --help              # Main help
./ddt glossary --help     # Glossary help
./ddt links --help        # Links help
./ddt version             # Show version
```

For detailed documentation, see [README.md](README.md).

