#!/bin/bash
set -e

# TOTP Manager - Test Coverage Report Generator
# Generates coverage reports for all packages with detailed metrics

echo "=== TOTP Manager Test Coverage ==="
echo ""

# Run tests with coverage
echo "Running tests with coverage..."
go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Display summary
echo ""
echo "=== Coverage Summary ==="
go tool cover -func=coverage.out | tail -n 1

# Check critical packages for 100% coverage
echo ""
echo "=== Critical Package Coverage (must be 100%) ==="
echo "Crypto package:"
go tool cover -func=coverage.out | grep "internal/crypto" || echo "No tests found"

echo ""
echo "TOTP package:"
go tool cover -func=coverage.out | grep "internal/totp" || echo "No tests found"

echo ""
echo "Storage package:"
go tool cover -func=coverage.out | grep "internal/storage" || echo "No tests found"

# Generate HTML report
echo ""
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report generated: coverage.html"

# Calculate overall coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo ""
echo "=== Coverage Check ==="
echo "Total coverage: ${COVERAGE}%"

# Check if coverage meets minimum requirement (80%)
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "❌ FAIL: Coverage ${COVERAGE}% is below minimum 80%"
    exit 1
else
    echo "✅ PASS: Coverage ${COVERAGE}% meets minimum 80% requirement"
fi
