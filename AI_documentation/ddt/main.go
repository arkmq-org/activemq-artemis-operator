package main

import (
	"fmt"
	"os"

	"github.com/arkmq-org/ddt/cmd"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "glossary":
		cmd.RunGlossary(os.Args[2:])
	case "links":
		cmd.RunLinks(os.Args[2:])
	case "version", "-v", "--version":
		fmt.Printf("ddt version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Developer Documentation Tools (ddt)

A unified CLI tool for maintaining developer documentation.

Usage:
  ddt <command> [flags]

Commands:
  glossary    Generate glossary from terms configuration
  links       Update documentation links to new commit
  version     Show version information
  help        Show this help message

Examples:
  # Generate glossary
  ddt glossary

  # Check documentation links (dry-run)
  ddt links

  # Update links to specific commit
  ddt links --update v1.2.3

  # Show help for specific command
  ddt glossary --help
  ddt links --help

For more information: https://github.com/arkmq-org/ddt`)
}
