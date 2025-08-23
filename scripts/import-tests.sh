#!/bin/bash

# Script to import test files from original ledger test suite

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
LEDGER_ROOT="$(dirname "$PROJECT_ROOT")/ledger"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Importing test files from original ledger...${NC}"

# Check if ledger directory exists
if [ ! -d "$LEDGER_ROOT" ]; then
    echo -e "${RED}Error: Ledger directory not found at $LEDGER_ROOT${NC}"
    exit 1
fi

# Create test fixture directories
mkdir -p "$PROJECT_ROOT/test/fixtures/baseline"
mkdir -p "$PROJECT_ROOT/test/fixtures/regress"
mkdir -p "$PROJECT_ROOT/test/fixtures/manual"
mkdir -p "$PROJECT_ROOT/test/fixtures/input"

# Import baseline tests
echo -e "${YELLOW}Importing baseline tests...${NC}"
if [ -d "$LEDGER_ROOT/test/baseline" ]; then
    cp -r "$LEDGER_ROOT/test/baseline"/*.test "$PROJECT_ROOT/test/fixtures/baseline/" 2>/dev/null || true
    cp -r "$LEDGER_ROOT/test/baseline"/*.dat "$PROJECT_ROOT/test/fixtures/baseline/" 2>/dev/null || true
    echo -e "${GREEN}✓ Imported baseline tests${NC}"
else
    echo -e "${YELLOW}⚠ Baseline test directory not found${NC}"
fi

# Import regression tests
echo -e "${YELLOW}Importing regression tests...${NC}"
if [ -d "$LEDGER_ROOT/test/regress" ]; then
    cp -r "$LEDGER_ROOT/test/regress"/*.test "$PROJECT_ROOT/test/fixtures/regress/" 2>/dev/null || true
    cp -r "$LEDGER_ROOT/test/regress"/*.dat "$PROJECT_ROOT/test/fixtures/regress/" 2>/dev/null || true
    echo -e "${GREEN}✓ Imported regression tests${NC}"
else
    echo -e "${YELLOW}⚠ Regression test directory not found${NC}"
fi

# Import manual tests
echo -e "${YELLOW}Importing manual tests...${NC}"
if [ -d "$LEDGER_ROOT/test/manual" ]; then
    cp -r "$LEDGER_ROOT/test/manual"/*.test "$PROJECT_ROOT/test/fixtures/manual/" 2>/dev/null || true
    echo -e "${GREEN}✓ Imported manual tests${NC}"
else
    echo -e "${YELLOW}⚠ Manual test directory not found${NC}"
fi

# Import input data files
echo -e "${YELLOW}Importing input data files...${NC}"
if [ -d "$LEDGER_ROOT/test/input" ]; then
    cp -r "$LEDGER_ROOT/test/input"/*.dat "$PROJECT_ROOT/test/fixtures/input/" 2>/dev/null || true
    cp -r "$LEDGER_ROOT/test/input"/*.ledger "$PROJECT_ROOT/test/fixtures/input/" 2>/dev/null || true
    echo -e "${GREEN}✓ Imported input data files${NC}"
else
    echo -e "${YELLOW}⚠ Input data directory not found${NC}"
fi

# Count imported files
BASELINE_COUNT=$(find "$PROJECT_ROOT/test/fixtures/baseline" -name "*.test" 2>/dev/null | wc -l)
REGRESS_COUNT=$(find "$PROJECT_ROOT/test/fixtures/regress" -name "*.test" 2>/dev/null | wc -l)
MANUAL_COUNT=$(find "$PROJECT_ROOT/test/fixtures/manual" -name "*.test" 2>/dev/null | wc -l)
INPUT_COUNT=$(find "$PROJECT_ROOT/test/fixtures/input" -name "*.dat" -o -name "*.ledger" 2>/dev/null | wc -l)

echo -e "${GREEN}Import complete!${NC}"
echo -e "  Baseline tests: $BASELINE_COUNT"
echo -e "  Regression tests: $REGRESS_COUNT"
echo -e "  Manual tests: $MANUAL_COUNT"
echo -e "  Input files: $INPUT_COUNT"
echo -e ""
echo -e "Total test files imported: $((BASELINE_COUNT + REGRESS_COUNT + MANUAL_COUNT))"