# Contributing to mcinit

Thank you for considering contributing to mcinit! This document provides guidelines for contributing.

## Development Setup

1. **Prerequisites**:
   - Go 1.21 or later
   - Git

2. **Clone the repository**:
   ```bash
   git clone https://github.com/jackh54/mcinit.git
   cd mcinit
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Build**:
   ```bash
   make build
   ```

5. **Run tests**:
   ```bash
   make test
   ```

## Project Structure

```
mcinit/
├── cmd/mcinit/          # Main application entry point
├── internal/            # Private application code
│   ├── cli/             # CLI commands
│   ├── config/          # Configuration management
│   ├── server/          # Server lifecycle management
│   ├── provider/        # Server jar providers (Vanilla, Paper, etc.)
│   ├── java/            # Java detection and validation
│   ├── cache/           # Cache and download management
│   ├── scripts/         # Startup script generation
│   ├── logs/            # Log reading and filtering
│   └── utils/           # Utility functions
├── pkg/                 # Public packages
│   └── jvmflags/        # JVM flags presets
└── scripts/             # Installation scripts
```

## Adding a New Server Provider

To add support for a new server type:

1. Create a new file in `internal/provider/` (e.g., `newserver.go`)
2. Implement the `Provider` interface
3. Register the provider in `internal/provider/registry.go`

## Testing

- Run all tests: `make test`
- Run with coverage: `make test-coverage`
- Run linter: `make lint`

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run tests and linting
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Code Style

- Follow Go best practices
- Use `gofmt` for formatting
- Write clear, concise comments
- Add tests for new features

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

