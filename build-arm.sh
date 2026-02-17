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
VERSION="1.0.0"
BUILD_DIR="build"
DIST_DIR="dist"

# Build configurations
GOOS="linux"
GOARCH="arm64"  # Change to arm for 32-bit ARM
CGO_ENABLED=0

echo "📦 Build Configuration:"
echo "   OS: $GOOS"
echo "   ARCH: $GOARCH"
echo "   CGO: $CGO_ENABLED"
echo ""

# Clean previous build
echo "🧹 Cleaning previous build..."
rm -rf $BUILD_DIR
rm -rf $DIST_DIR
mkdir -p $BUILD_DIR
mkdir -p $DIST_DIR

# Step 1: Tidy dependencies
echo "📚 Tidying Go modules..."
go mod tidy

# Step 2: Build binary
echo "🔨 Building binary..."
GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED go build \
    -ldflags="-s -w -X main.Version=$VERSION" \
    -o $BUILD_DIR/$APP_NAME \
    cmd/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Binary built successfully!${NC}"
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

# Step 3: Copy necessary files
echo "📋 Copying necessary files..."

# Copy web directory
cp -r web $BUILD_DIR/

# Copy example env
cp .env.example $BUILD_DIR/.env.example

# Copy migrations if exists
if [ -d "migrations" ]; then
    cp -r migrations $BUILD_DIR/
fi

# Step 4: Create package
echo "📦 Creating deployment package..."
cd $BUILD_DIR

# Create tarball
tar -czf ../$DIST_DIR/${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz .

cd ..

# Step 5: Create installer package (single binary with embedded files)
echo "🎁 Creating single binary installer..."

# For now, the installer will copy files from web directory
# In production, we can use go:embed to embed files into binary

# Make the binary executable
chmod +x $BUILD_DIR/$APP_NAME

# Copy to dist with installer suffix
cp $BUILD_DIR/$APP_NAME $DIST_DIR/${APP_NAME}-installer-${GOOS}-${GOARCH}
chmod +x $DIST_DIR/${APP_NAME}-installer-${GOOS}-${GOARCH}

# Step 6: Show file sizes
echo ""
echo "📊 Build Results:"
echo "   Binary size: $(du -h $BUILD_DIR/$APP_NAME | cut -f1)"
echo "   Package size: $(du -h $DIST_DIR/${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz | cut -f1)"
echo ""

# Step 7: Show checksums
echo "🔐 Checksums:"
if command -v sha256sum &> /dev/null; then
    sha256sum $DIST_DIR/${APP_NAME}-installer-${GOOS}-${GOARCH}
elif command -v shasum &> /dev/null; then
    shasum -a 256 $DIST_DIR/${APP_NAME}-installer-${GOOS}-${GOARCH}
fi
echo ""

# Step 8: Create README for deployment
cat > $DIST_DIR/README.md << 'EOF'
# Amaliah Ramadhan - Deployment Guide

## Installation on Armbian/ARM Linux

### Prerequisites
- Armbian or any ARM-based Linux distribution
- Root or sudo access
- Systemd installed

### Quick Installation

1. **Transfer the installer to your server:**
   ```bash
   scp amaliah-ramadhan-installer-linux-arm64 user@your-server:/tmp/
   ```

2. **SSH to your server:**
   ```bash
   ssh user@your-server
   ```

3. **Run the installer:**
   ```bash
   cd /tmp
   chmod +x amaliah-ramadhan-installer-linux-arm64
   sudo ./amaliah-ramadhan-installer-linux-arm64 -install
   ```

4. **Follow the wizard:**
   - Choose "1" for New Installation
   - Enter desired port (default: 8080)
   - Wait for installation to complete

### Manual Installation (from tarball)

1. **Extract the package:**
   ```bash
   tar -xzf amaliah-ramadhan-*.tar.gz
   cd amaliah-ramadhan
   ```

2. **Create installation directory:**
   ```bash
   sudo mkdir -p /opt/amaliah-ramadhan
   sudo cp -r * /opt/amaliah-ramadhan/
   ```

3. **Create config file:**
   ```bash
   sudo cp /opt/amaliah-ramadhan/.env.example /opt/amaliah-ramadhan/.env
   sudo nano /opt/amaliah-ramadhan/.env  # Edit as needed
   ```

4. **Create systemd service:**
   ```bash
   sudo nano /etc/systemd/system/amaliah-ramadhan.service
   ```
   
   Paste this content:
   ```ini
   [Unit]
   Description=Amaliah Ramadhan - Monitoring Ibadah Harian
   After=network.target

   [Service]
   Type=simple
   User=root
   WorkingDirectory=/opt/amaliah-ramadhan
   ExecStart=/opt/amaliah-ramadhan/amaliah-ramadhan
   Restart=always
   RestartSec=5
   StandardOutput=append:/var/log/amaliah-ramadhan.log
   StandardError=append:/var/log/amaliah-ramadhan.error.log

   [Install]
   WantedBy=multi-user.target
   ```

5. **Enable and start service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable amaliah-ramadhan
   sudo systemctl start amaliah-ramadhan
   ```

6. **Check status:**
   ```bash
   sudo systemctl status amaliah-ramadhan
   ```

### Update Existing Installation

Run installer and choose option "2" for Update:
```bash
sudo ./amaliah-ramadhan-installer-linux-arm64 -install
```

### Access the Application

Open your browser and navigate to:
```
http://your-server-ip:8080
```

Default login:
- Username: `admin`
- Password: `admin123`

**⚠️ Change the default password after first login!**

### Useful Commands

```bash
# Check status
sudo systemctl status amaliah-ramadhan

# View logs
sudo journalctl -u amaliah-ramadhan -f

# Restart service
sudo systemctl restart amaliah-ramadhan

# Stop service
sudo systemctl stop amaliah-ramadhan

# Start service
sudo systemctl start amaliah-ramadhan
```

### Troubleshooting

**Service won't start:**
1. Check logs: `sudo journalctl -u amaliah-ramadhan -n 50`
2. Check permissions: `ls -la /opt/amaliah-ramadhan`
3. Check config: `cat /opt/amaliah-ramadhan/.env`

**Can't access from browser:**
1. Check if service is running: `sudo systemctl status amaliah-ramadhan`
2. Check firewall: `sudo ufw status`
3. Check port: `sudo netstat -tlnp | grep 8080`

**Database errors:**
1. Check database file exists: `ls -la /opt/amaliah-ramadhan/amaliah.db`
2. Check permissions: `sudo chown -R root:root /opt/amaliah-ramadhan`

### Support

For issues and support, please contact the administrator.
EOF

echo -e "${GREEN}✅ Build completed successfully!${NC}"
echo ""
echo "📦 Output files:"
echo "   • Installer: $DIST_DIR/${APP_NAME}-installer-${GOOS}-${GOARCH}"
echo "   • Package:   $DIST_DIR/${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz"
echo "   • README:    $DIST_DIR/README.md"
echo ""
echo "🚀 To install on Armbian server:"
echo "   1. Copy installer to server"
echo "   2. Run: sudo ./${APP_NAME}-installer-${GOOS}-${GOARCH} -install"
echo ""
echo -e "${YELLOW}📖 Check $DIST_DIR/README.md for detailed instructions${NC}"
echo ""
