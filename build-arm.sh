#!/bin/bash

# Build script for Amaliah Ramadhan - ARM Binary
# For Armbian/ARM architecture deployment

set -e

echo "╔══════════════════════════════════════════════════════════╗"
echo "║     Building Amaliah Ramadhan for ARM Architecture      ║"
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
GOARCH="arm64"  # Change to arm for 32-bit ARM
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

# Step 1: Install Dependencies & Build JS/CSS
echo "📦 Installing front-end dependencies..."
npm install

echo "🎨 Building CSS..."
npm run build:css

# Step 2: Tidy dependencies
echo "📚 Tidying Go modules..."
go mod tidy

# Step 3: Build binary
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

# Step 4: Verify binary
ls -lh $BUILD_DIR/$APP_NAME

# Step 5: Move to dist
echo "📦 Moving binary to dist..."
OUTPUT_NAME="${APP_NAME}-${VERSION}-${DATE}-${GOOS}-${GOARCH}"
cp $BUILD_DIR/$APP_NAME $DIST_DIR/$OUTPUT_NAME
chmod +x $DIST_DIR/$OUTPUT_NAME

# Step 6: Create README
cat > $DIST_DIR/README.md << 'EOF'
# Amaliah Ramadhan - Single Binary Deployment (ARM)

## Installation on Armbian/ARM Linux

1. **Upload** the binary to your server:
   ```bash
   scp amaliah-ramadhan-*-linux-arm64 user@server:/tmp/
   ```

2. **Run Installer**:
   ```bash
   ssh user@server
   chmod +x /tmp/amaliah-ramadhan-*-linux-arm64
   sudo /tmp/amaliah-ramadhan-*-linux-arm64 -install
   ```

3. **Follow the Wizard**.
   - Choose "1" for New Installation.
   - The application will handle everything.

### Update Existing Installation

1. Upload the new binary.
2. Run it with `-install`.
3. Choose "2" for Update.
EOF

echo ""
echo "📊 Build Results:"
echo "   Binary: $DIST_DIR/$OUTPUT_NAME"
echo "   Size:   $(du -h $DIST_DIR/$OUTPUT_NAME | cut -f1)"
echo ""
echo "✅ Build completed!"
echo "🚀 To install:"
echo "   1. Copy $OUTPUT_NAME to server"
echo "   2. Run: sudo ./$OUTPUT_NAME -install"
echo ""
