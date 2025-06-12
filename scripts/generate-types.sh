#!/bin/bash
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Generating TypeScript types...${NC}"

# Change to desktop directory and run the npm script
cd desktop
if npm run generate-types; then
    echo -e "${GREEN}✓ TypeScript types generated successfully${NC}"
    exit 0
else
    echo -e "${RED}✗ Failed to generate TypeScript types${NC}"
    exit 1
fi