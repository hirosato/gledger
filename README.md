# Gledger - Go Implementation of Ledger-CLI

Gledger is a Go implementation of the popular [ledger-cli](https://www.ledger-cli.org/) double-entry accounting system. It aims to be fully compatible with the original ledger while providing a modern, maintainable codebase using Hexagonal Architecture principles.

## Project Status

🚧 **Under Development** - This project is in early development stages.

## Architecture

The project follows Hexagonal Architecture (Ports and Adapters) pattern:

```
gledger/
├── domain/         # Core business logic (no external dependencies)
├── application/    # Use cases and application services
├── infrastructure/ # External dependencies (file I/O, parsers)
├── interfaces/     # Input/output adapters (CLI)
├── cmd/           # Application entry point
└── test/          # Test suites and fixtures
```

## Building

### Prerequisites

- Go 1.21 or higher
- Make
- ledger-cli (for comparison tests)

### Build Commands

```bash
# Build the binary
make build

# Run all tests
make test

# Run specific test types
make test-unit        # Unit tests only
make test-integration # Integration tests only
make test-compare     # Comparison tests with ledger-cli

# Format and lint
make fmt
make lint

# Clean build artifacts
make clean

# See all available commands
make help
```

## Testing Strategy

Gledger maintains spec equivalence with the original ledger-cli through:

1. **Imported Test Suite**: All test cases from the original ledger are imported and must pass
2. **Comparison Testing**: Output is compared byte-for-byte with ledger-cli
3. **Go-Idiomatic Tests**: Additional unit and integration tests following Go best practices

### Running Tests

```bash
# Import test files from original ledger
make import-tests

# Run all tests
make test

# Run comparison tests
make test-compare

# Run with coverage
make coverage
```

## Development Progress

See [TASKS.md](../TASKS.md) for detailed implementation plan and progress tracking.

### Current Phase: Project Setup ✅

- [x] Initialize Go module
- [x] Set up project directory structure
- [x] Create Makefile
- [x] Set up testing framework
- [x] Configure CI/CD
- [x] Create test harness

### Next Phase: Core Domain Models

- [ ] Implement Account entity
- [ ] Implement Transaction entity
- [ ] Implement Posting entity
- [ ] Implement Amount value object
- [ ] Implement Balance aggregate

## Compatibility Goals

- 100% compatibility with ledger-cli file format
- Identical output for all commands
- Pass all original ledger test cases
- Performance within 2x of original C++ implementation

## Contributing

This project is under active development. Contributions are welcome once the core architecture is established.

## License

[To be determined - will match original ledger-cli license]

## Acknowledgments

This project is based on the excellent work of the [ledger-cli](https://www.ledger-cli.org/) project and aims to provide a compatible implementation in Go.