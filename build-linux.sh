#!/bin/bash

# Build script for Amaliah Ramadhan - Linux AMD64 Binary
# For Standard Linux Server (x86_64) deployment

set -e

echo "╔══════════════════════════════════════════════════════════╗"
echo "║     Building Amaliah Ramadhan for Linux (AMD64)         ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables
APP_NAME="amaliah-ramadhan"
# Load version from .env if available
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
    VERSION=${APP_VERSION:-"1.0.0"}
else
    VERSION="1.0.0"
fi
DATE=$(date +%Y%m%d)
BUILD_DIR="build"
DIST_DIR="dist"

# Build configurations
GOOS="linux"
GOARCH="amd64"
CGO_ENABLED=0

echo "📦 Build Configuration:"
echo "   OS: $GOOS"
echo "   ARCH: $GOARCH"
echo "   Date: $DATE"
echo "   CGO: $CGO_ENABLED"
echo ""

# Clean previous build
echo "🧹 Cleaning previous build..."
rm -rf $BUILD_DIR
rm -rf $DIST_DIR
mkdir -p $BUILD_DIR
mkdir -p $DIST_DIR

# Step 0: Build Frontend Assets
echo "🎨 Building frontend assets..."
if command -v npm &> /dev/null; then
    echo "   Running npm install..."
    npm install
    echo "   Compiling CSS..."
    npm run build:css
else
    echo -e "${YELLOW}⚠️  npm not found! Skipping frontend build.${NC}"
    echo -e "${YELLOW}   Ensure web/static/css/output.css exists or install Node.js.${NC}"
fi

# Step 1: Tidy dependencies
echo "📚 Tidying Go modules..."
go mod tidy

# Step 2: Build binary
echo "🔨 Building binary..."
GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED go build \
    -ldflags="-s -w -X main.Version=$VERSION -X main.BuildDate=$DATE" \
    -o $BUILD_DIR/$APP_NAME \
    cmd/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Binary built successfully!${NC}"
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

# Step 3: Verify binary
ls -lh $BUILD_DIR/$APP_NAME

# Step 4: Move to dist
echo "📦 Moving binary to dist..."
OUTPUT_NAME="${APP_NAME}-${VERSION}-${DATE}-${GOOS}-${GOARCH}"
cp $BUILD_DIR/$APP_NAME $DIST_DIR/$OUTPUT_NAME
chmod +x $DIST_DIR/$OUTPUT_NAME

# Create README
cat > $DIST_DIR/README.md << 'EOF'
# Amaliah Ramadhan - Single Binary Deployment

## Installation

1. **Upload** the binary to your server:
   ```bash
   scp amaliah-ramadhan-*-linux-amd64 user@server:/tmp/
   ```

2. **Run Installer**:
   ```bash
   ssh user@server
   chmod +x /tmp/amaliah-ramadhan-*-linux-amd64
   sudo /tmp/amaliah-ramadhan-*-linux-amd64 -install
   ```

3. **Follow the Wizard**.
EOF

echo ""
echo "📊 Build Results:"
echo "   Binary: $DIST_DIR/$OUTPUT_NAME"
echo "   Size:   $(du -h $DIST_DIR/$OUTPUT_NAME | cut -f1)"
echo ""
echo "✅ Build completed!"
