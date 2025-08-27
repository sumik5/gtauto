// Package main provides a tool to automatically create git tags with CHANGELOG content
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var version = "1.0.0" // Set during build

const (
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
	colorReset  = "\033[0m"
)

func main() {
	tagName := flag.String("tag", "", "Tag name to create (required)")
	changelogFile := flag.String("changelog", "CHANGELOG.md", "Path to CHANGELOG file")
	showHelp := flag.Bool("h", false, "Show help message")
	showHelpLong := flag.Bool("help", false, "Show help message")
	showVersion := flag.Bool("version", false, "Show version information")
	force := flag.Bool("force", false, "Force overwrite existing tag without confirmation")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gtauto: Git tag automation with CHANGELOG support\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  gtauto --tag <tag_name> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  gtauto --tag v1.0.0\n")
		fmt.Fprintf(os.Stderr, "  gtauto --tag v1.0.0 --changelog path/to/CHANGELOG.md\n")
		fmt.Fprintf(os.Stderr, "  gtauto --tag v1.0.0 --force\n")
	}

	flag.Parse()

	if *showHelp || *showHelpLong {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("gtauto version %s\n", version)
		os.Exit(0)
	}

	if *tagName == "" {
		printError("--tag option is required")
		flag.Usage()
		os.Exit(1)
	}

	// Check if we're in a git repository
	if err := checkGitRepository(); err != nil {
		printError(fmt.Sprintf("Not a git repository: %v", err))
		os.Exit(1)
	}

	// Check if CHANGELOG file exists
	if _, err := os.Stat(*changelogFile); os.IsNotExist(err) {
		printError(fmt.Sprintf("CHANGELOG file not found: %s", *changelogFile))
		os.Exit(1)
	}

	// Check if tag already exists
	if tagExists(*tagName) {
		if !*force {
			printWarning(fmt.Sprintf("Tag '%s' already exists", *tagName))
			if !confirmOverwrite() {
				fmt.Println("Operation cancelled")
				os.Exit(0)
			}
		}
		// Delete existing tag
		if err := deleteTag(*tagName); err != nil {
			printError(fmt.Sprintf("Failed to delete existing tag: %v", err))
			os.Exit(1)
		}
	}

	printSuccess(fmt.Sprintf("Extracting CHANGELOG entry for '%s'...", *tagName))

	// Extract changelog entry
	changelogEntry, err := extractChangelogEntry(*tagName, *changelogFile)
	if err != nil {
		printWarning(fmt.Sprintf("Could not find CHANGELOG entry for '%s'", *tagName))
		changelogEntry = fmt.Sprintf("Release %s", *tagName)
	} else {
		printSuccess("Found CHANGELOG entry")
	}

	// Create annotated tag
	printSuccess(fmt.Sprintf("Creating tag '%s'...", *tagName))
	fmt.Println("\nTag message:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println(changelogEntry)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println()

	if err := createTag(*tagName, changelogEntry); err != nil {
		printError(fmt.Sprintf("Failed to create tag: %v", err))
		os.Exit(1)
	}

	printSuccess(fmt.Sprintf("✓ Tag '%s' created successfully", *tagName))
	fmt.Println("\nTo push this tag to remote:")
	fmt.Printf("  git push origin %s\n", *tagName)
	fmt.Println("\nTo push all tags:")
	fmt.Println("  git push --tags")
}

func checkGitRepository() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run()
}

func tagExists(tagName string) bool {
	cmd := exec.Command("git", "tag", "-l", tagName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == tagName
}

func deleteTag(tagName string) error {
	cmd := exec.Command("git", "tag", "-d", tagName)
	return cmd.Run()
}

func confirmOverwrite() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to overwrite it? (y/N): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func extractChangelogEntry(tagName, changelogFile string) (string, error) {
	file, err := os.Open(changelogFile)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	// Remove 'v' prefix if present to match version number
	version := strings.TrimPrefix(tagName, "v")

	// Pattern to match version headers like ## [v1.0.0] or ## v1.0.0
	versionPattern := fmt.Sprintf(`^##\s+\[?v?%s\]?`, regexp.QuoteMeta(version))
	versionRegex := regexp.MustCompile(versionPattern)
	nextVersionRegex := regexp.MustCompile(`^##\s+\[?v?[0-9]+\.[0-9]+`)

	scanner := bufio.NewScanner(file)
	var inSection bool
	var content strings.Builder
	var sectionFound bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is the version we're looking for
		if versionRegex.MatchString(line) {
			inSection = true
			sectionFound = true
			content.WriteString(line)
			content.WriteString("\n")
			continue
		}

		// Check if we've reached the next version section
		if inSection && nextVersionRegex.MatchString(line) {
			break
		}

		// If we're in the right section, collect the content
		if inSection {
			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if !sectionFound {
		return "", fmt.Errorf("version %s not found in changelog", tagName)
	}

	// Trim trailing empty lines
	result := strings.TrimRight(content.String(), "\n")
	return result, nil
}

func createTag(tagName, message string) error {
	cmd := exec.Command("git", "tag", "-a", tagName, "-m", message)
	return cmd.Run()
}

func printError(message string) {
	fmt.Printf("%sError: %s%s\n", colorRed, message, colorReset)
}

func printWarning(message string) {
	fmt.Printf("%sWarning: %s%s\n", colorYellow, message, colorReset)
}

func printSuccess(message string) {
	fmt.Printf("%s%s%s\n", colorGreen, message, colorReset)
}
