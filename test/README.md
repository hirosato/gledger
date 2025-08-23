# Test Suite

This directory contains all test files and test infrastructure.

## Structure

- **fixtures/**: Test data files from original ledger suite
- **specs/**: Test specifications extracted from original tests
- **integration/**: End-to-end tests comparing with ledger-cli
- **unit/**: Unit tests for individual components

## Test Philosophy

- Maintain spec equivalence with original ledger-cli
- Test files are imported from original ledger test suite
- Automated comparison testing with ledger-cli
- Go-idiomatic testing where appropriate

## Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run comparison tests with ledger-cli
make test-compare

# Run specific test suite
go test ./test/unit/...
```