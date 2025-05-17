#!/bin/bash

echo "MCP-Pandoc-Go Installer for Linux/Ubuntu"
echo "========================================"

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if running with sudo privileges
if [ "$EUID" -ne 0 ]; then
    echo "ERROR: This script requires administrator privileges."
    echo "Please run with sudo: sudo ./install_linux.sh"
    exit 1
fi

# Check for required tools
if ! command_exists git; then
    echo "ERROR: Git is not installed."
    echo "Please install Git using: sudo apt-get install git"
    exit 1
fi

if ! command_exists go; then
    echo "ERROR: Go is not installed."
    echo "Please install Go using: sudo apt-get install golang-go"
    echo "Or download from: https://golang.org/dl/"
    exit 1
fi

if ! command_exists pandoc; then
    echo "WARNING: Pandoc is not installed or not in PATH."
    echo "Please install Pandoc using: sudo apt-get install pandoc"
    echo "Or download from: https://pandoc.org/installing.html"
    echo "After installation, run this script again."
    exit 1
fi

# Get actual user (not sudo)
ACTUAL_USER=$(logname || echo $SUDO_USER)
if [ -z "$ACTUAL_USER" ]; then
    echo "WARNING: Could not determine the actual user. Some steps might fail."
    ACTUAL_USER=$USER
fi

# Create installation directory
echo "Creating installation directory..."
mkdir -p /opt/SnowWhiteAI/MCP_servers/mcp-pandoc-go
if [ $? -ne 0 ]; then
    echo "Failed to create directories."
    exit 1
fi

# Clone the repository
echo "Cloning repository..."
cd /opt/SnowWhiteAI/MCP_servers
if [ -d "mcp-pandoc-go/.git" ]; then
    echo "Repository already exists, updating..."
    cd mcp-pandoc-go
    git pull
else
    echo "Cloning fresh repository..."
    git clone https://github.com/otinoff/mcp-pandoc-go.git
    cd mcp-pandoc-go
fi

# Build the application
echo "Building MCP-Pandoc-Go..."
go build -o pandoc-mcp-go ./cmd/main.go

# Set proper ownership and permissions
chown -R $ACTUAL_USER:$ACTUAL_USER /opt/SnowWhiteAI
chmod +x /opt/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go

# Create cursor directory if it doesn't exist
echo "Setting up Cursor configuration..."
CURSOR_CONFIG_DIR="/home/$ACTUAL_USER/.cursor/mcp/user"
mkdir -p $CURSOR_CONFIG_DIR
chown -R $ACTUAL_USER:$ACTUAL_USER /home/$ACTUAL_USER/.cursor

# Create or update MCP configuration
cat > $CURSOR_CONFIG_DIR/config.json << EOF
{
  "servers": {
    "pandoc_mcp_go": {
      "type": "stdio",
      "command": "/opt/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go",
      "cwd": "/opt/SnowWhiteAI/MCP_servers/mcp-pandoc-go"
    }
  },
  "roots": {
    "pandoc": {
      "type": "document-converter",
      "server": "pandoc_mcp_go"
    }
  }
}
EOF

chown $ACTUAL_USER:$ACTUAL_USER $CURSOR_CONFIG_DIR/config.json

# Create symbolic link to make it accessible from PATH
echo "Creating symbolic link..."
ln -sf /opt/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go /usr/local/bin/pandoc-mcp-go

echo ""
echo "Installation completed successfully!"
echo "MCP-Pandoc-Go has been installed to: /opt/SnowWhiteAI/MCP_servers/mcp-pandoc-go"
echo "Configuration has been set up for Cursor IDE."
echo ""
echo "Please restart Cursor if it's already running."
echo "" 