# gtauto

Git tag automation tool with CHANGELOG support.

## Features

- Automatically create git tags with CHANGELOG content as annotation
- Extract specific version entries from CHANGELOG.md
- Support Keep a Changelog format
- Color output for better visibility
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Using mise

```bash
mise run install
```

### Using make

```bash
make install
```

### Manual installation

```bash
go build -o gtauto main.go
sudo mv gtauto /usr/local/bin/
```

### Download pre-built binaries

Download the latest release from [GitHub Releases](https://github.com/shivase/gtauto/releases).

## Usage

### Basic usage

```bash
gtauto --tag v1.0.0
```

This will:
1. Extract the v1.0.0 section from CHANGELOG.md
2. Create an annotated git tag with the extracted content

### Options

```bash
gtauto --tag <tag_name> [options]

Options:
  --tag <tag_name>        Tag name to create (required)
  --changelog <file>      Path to CHANGELOG file (default: CHANGELOG.md)
  --force                 Force overwrite existing tag without confirmation
  --version              Show version information
  --help                 Show help message
```

### Examples

```bash
# Create a tag with default CHANGELOG.md
gtauto --tag v1.0.0

# Use a different changelog file
gtauto --tag v1.0.0 --changelog docs/CHANGELOG.md

# Force overwrite existing tag
gtauto --tag v1.0.0 --force

# Show version
gtauto --version
```

## CHANGELOG Format

`gtauto` expects the CHANGELOG to follow the [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
# Changelog

## [v1.0.1] - 2025-08-27

### Added
- New feature A
- New feature B

### Fixed
- Bug fix 1

## [v1.0.0] - 2025-08-26

### Added
- Initial release
```

The tool will extract the entire section for the specified version, including all subsections (Added, Changed, Fixed, etc.).

## Development

### Prerequisites

- Go 1.21 or higher
- mise (optional, for task management)

### Setup

```bash
# Clone the repository
git clone https://github.com/shivase/gtauto.git
cd gtauto

# Install dependencies
go mod download
```

### Running tests

```bash
# Using mise
mise run test

# Using make
make test

# Using go directly
go test -v ./...
```

### Building

```bash
# Using mise
mise run build

# Using make
make build

# Using go directly
go build -o build/gtauto main.go
```

### Development workflow

```bash
# Format code
mise run fmt

# Run linter
mise run lint

# Run all development tasks
mise run dev
```

### Building for all platforms

```bash
# Using mise
mise run release

# Using make
make build-all
```

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **Lint**: Code formatting and static analysis
- **Test**: Unit tests
- **Build**: Cross-platform binary compilation
- **Release**: Automatic release creation on tag push

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

Inspired by [changelog-update](https://github.com/shivase/changelog-update) project structure and workflow.