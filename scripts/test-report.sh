#!/bin/bash

# test-report.sh - Standard Go testing and reporting
# Uses standard Go tools + industry standard reporting tools

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORTS_DIR="${PROJECT_ROOT}/reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo -e "${CYAN}🛡️  Praetorian Go - Standard Test Report${NC}"
echo -e "${CYAN}=====================================${NC}"

# Create reports directory
mkdir -p "$REPORTS_DIR"

cd "$PROJECT_ROOT"

echo -e "${CYAN}🧪 Running tests with coverage...${NC}"

# 1. Run tests with coverage (standard Go)
go test -v -coverprofile="${REPORTS_DIR}/coverage.out" -covermode=atomic ./...

echo -e "${GREEN}✅ Tests completed${NC}"

# 2. Generate HTML coverage report (standard Go)
echo -e "${CYAN}📊 Generating HTML coverage report...${NC}"
go tool cover -html="${REPORTS_DIR}/coverage.out" -o "${REPORTS_DIR}/coverage.html"

echo -e "${GREEN}✅ HTML report: ${REPORTS_DIR}/coverage.html${NC}"

# 3. Generate functional coverage report (standard Go)
echo -e "${CYAN}📋 Generating functional coverage report...${NC}"
go tool cover -func="${REPORTS_DIR}/coverage.out" > "${REPORTS_DIR}/coverage-func.txt"

# Extract total coverage
TOTAL_COVERAGE=$(tail -n 1 "${REPORTS_DIR}/coverage-func.txt" | awk '{print $NF}')
echo -e "${GREEN}✅ Total Coverage: ${TOTAL_COVERAGE}${NC}"

# 4. Generate JUnit XML report (industry standard)
echo -e "${CYAN}📄 Generating JUnit XML report...${NC}"
go test -v ./... 2>&1 | go-junit-report > "${REPORTS_DIR}/junit-${TIMESTAMP}.xml"

echo -e "${GREEN}✅ JUnit report: ${REPORTS_DIR}/junit-${TIMESTAMP}.xml${NC}"

# 5. Generate JSON coverage report (standard Go + gocov)
if command -v gocov >/dev/null 2>&1; then
    echo -e "${CYAN}📄 Generating JSON coverage report...${NC}"
    gocov test ./... | gocov-xml > "${REPORTS_DIR}/coverage-${TIMESTAMP}.xml"
    echo -e "${GREEN}✅ XML coverage report: ${REPORTS_DIR}/coverage-${TIMESTAMP}.xml${NC}"
else
    echo -e "${YELLOW}⚠️  gocov not installed, skipping JSON report${NC}"
fi

# 6. Summary
echo ""
echo -e "${GREEN}🎉 Test reporting completed!${NC}"
echo ""
echo -e "${CYAN}📁 Reports generated:${NC}"
echo -e "  📊 HTML Coverage: ${REPORTS_DIR}/coverage.html"
echo -e "  📋 Functional: ${REPORTS_DIR}/coverage-func.txt"
echo -e "  📄 JUnit XML: ${REPORTS_DIR}/junit-${TIMESTAMP}.xml"
echo -e "  📈 Total Coverage: ${TOTAL_COVERAGE}"
echo ""

# Open HTML report if on macOS
if [[ "$OSTYPE" == "darwin"* ]] && command -v open >/dev/null 2>&1; then
    echo -e "${CYAN}🌐 Opening HTML report...${NC}"
    open "${REPORTS_DIR}/coverage.html"
fi
