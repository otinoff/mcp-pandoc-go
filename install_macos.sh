#!/bin/bash

echo "MCP-Pandoc-Go Installer for macOS"
echo "================================="

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check for required tools
if ! command_exists git; then
    echo "ERROR: Git is not installed."
    echo "Please install Git using Homebrew: brew install git"
    echo "Or download from: https://git-scm.com/download/mac"
    exit 1
fi

if ! command_exists go; then
    echo "ERROR: Go is not installed."
    echo "Please install Go using Homebrew: brew install go"
    echo "Or download from: https://golang.org/dl/"
    exit 1
fi

if ! command_exists pandoc; then
    echo "WARNING: Pandoc is not installed or not in PATH."
    echo "Please install Pandoc using Homebrew: brew install pandoc"
    echo "Or download from: https://pandoc.org/installing.html"
    echo "After installation, run this script again."
    exit 1
fi

# Create installation directory
echo "Creating installation directory..."
sudo mkdir -p /Applications/SnowWhiteAI/MCP_servers/mcp-pandoc-go
if [ $? -ne 0 ]; then
    echo "Failed to create directories. Do you have administrator privileges?"
    exit 1
fi

# Set proper permissions
sudo chown -R $USER:staff /Applications/SnowWhiteAI

# Clone the repository
echo "Cloning repository..."
cd /Applications/SnowWhiteAI/MCP_servers
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

# Create cursor directory if it doesn't exist
echo "Setting up Cursor configuration..."
mkdir -p ~/.cursor/mcp/user

# Create or update MCP configuration
cat > ~/.cursor/mcp/user/config.json << EOF
{
  "servers": {
    "pandoc_mcp_go": {
      "type": "stdio",
      "command": "/Applications/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go",
      "cwd": "/Applications/SnowWhiteAI/MCP_servers/mcp-pandoc-go"
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

# Create symbolic link to make it accessible from PATH
echo "Creating symbolic link..."
sudo ln -sf /Applications/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go /usr/local/bin/pandoc-mcp-go

echo ""
echo "Installation completed successfully!"
echo "MCP-Pandoc-Go has been installed to: /Applications/SnowWhiteAI/MCP_servers/mcp-pandoc-go"
echo "Configuration has been set up for Cursor IDE."
echo ""
echo "Please restart Cursor if it's already running."
echo "" 