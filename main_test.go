package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractChangelogEntry(t *testing.T) {
	tests := []struct {
		name             string
		tagName          string
		changelogContent string
		wantContent      string
		wantErr          bool
	}{
		{
			name:    "extract version with brackets",
			tagName: "v1.0.1",
			changelogContent: `# Changelog

## [v1.0.1] - 2025-08-27

### Added
- New feature A
- New feature B

### Fixed
- Bug fix 1

## [v1.0.0] - 2025-08-26

### Added
- Initial release`,
			wantContent: `## [v1.0.1] - 2025-08-27

### Added
- New feature A
- New feature B

### Fixed
- Bug fix 1`,
			wantErr: false,
		},
		{
			name:    "extract version without brackets",
			tagName: "v2.0.0",
			changelogContent: `# Changelog

## v2.0.0 - 2025-08-27

### Added
- Major feature

## v1.0.0 - 2025-08-26

### Added
- Initial release`,
			wantContent: `## v2.0.0 - 2025-08-27

### Added
- Major feature`,
			wantErr: false,
		},
		{
			name:    "version not found",
			tagName: "v3.0.0",
			changelogContent: `# Changelog

## [v1.0.0] - 2025-08-26

### Added
- Initial release`,
			wantContent: "",
			wantErr:     true,
		},
		{
			name:    "handle version with and without v prefix",
			tagName: "1.0.0",
			changelogContent: `# Changelog

## [v1.0.0] - 2025-08-26

### Added
- Initial release`,
			wantContent: `## [v1.0.0] - 2025-08-26

### Added
- Initial release`,
			wantErr: false,
		},
		{
			name:    "extract middle version",
			tagName: "v1.0.1",
			changelogContent: `# Changelog

## [v1.0.2] - 2025-08-28

### Added
- Latest feature

## [v1.0.1] - 2025-08-27

### Added
- Middle feature

### Fixed
- Middle bug

## [v1.0.0] - 2025-08-26

### Added
- Initial release`,
			wantContent: `## [v1.0.1] - 2025-08-27

### Added
- Middle feature

### Fixed
- Middle bug`,
			wantErr: false,
		},
		{
			name:    "handle trailing newlines",
			tagName: "v1.0.0",
			changelogContent: `# Changelog

## [v1.0.0] - 2025-08-26

### Added
- Initial release


`,
			wantContent: `## [v1.0.0] - 2025-08-26

### Added
- Initial release`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary changelog file
			tmpDir := t.TempDir()
			changelogFile := filepath.Join(tmpDir, "CHANGELOG.md")

			err := os.WriteFile(changelogFile, []byte(tt.changelogContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test changelog: %v", err)
			}

			// Test the extraction
			got, err := extractChangelogEntry(tt.tagName, changelogFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("extractChangelogEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Normalize whitespace for comparison
				gotNormalized := strings.TrimSpace(got)
				wantNormalized := strings.TrimSpace(tt.wantContent)

				if gotNormalized != wantNormalized {
					t.Errorf("extractChangelogEntry() content mismatch\nGot:\n%s\n\nWant:\n%s", got, tt.wantContent)
				}
			}
		})
	}
}

func TestTagExists(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		expected bool
		skip     bool // Skip tests that require actual git repo
	}{
		{
			name:     "non-existent tag",
			tagName:  "v999.999.999",
			expected: false,
			skip:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping test that requires git repository")
			}

			// This test will work in CI/CD as it doesn't require an existing tag
			got := tagExists(tt.tagName)
			if got != tt.expected {
				t.Errorf("tagExists() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckGitRepository(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "non-git directory",
			path:    t.TempDir(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to test directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer func() {
				_ = os.Chdir(originalDir)
			}()

			if err := os.Chdir(tt.path); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}

			err = checkGitRepository()
			if (err != nil) != tt.wantErr {
				t.Errorf("checkGitRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfirmOverwrite(t *testing.T) {
	// This test would require mocking stdin, which is complex
	// For now, we'll skip interactive tests
	t.Skip("Skipping interactive test")
}

func TestColorOutput(t *testing.T) {
	// Test that color constants are defined correctly
	tests := []struct {
		name  string
		color string
		want  string
	}{
		{"red color", colorRed, "\033[0;31m"},
		{"green color", colorGreen, "\033[0;32m"},
		{"yellow color", colorYellow, "\033[1;33m"},
		{"reset color", colorReset, "\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.want {
				t.Errorf("Color constant = %q, want %q", tt.color, tt.want)
			}
		})
	}
}
