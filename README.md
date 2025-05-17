# MCP-Pandoc

A simple and powerful document conversion server for Cursor IDE using Pandoc.

## When to use this

✅ **When you need to convert documents between different formats**  
Converting markdown to Word, HTML to PDF, or any combination of supported formats.

✅ **When you want to generate professional documents in Cursor**  
Creating Word documents, PDFs, or other formats directly from your markdown content.

✅ **When you need consistent document branding**  
Automatically adds copyright and branding to all generated documents.

✅ **When you need a reliable, high-quality document converter**  
Built on the industry-standard Pandoc conversion engine.

## Features

- Fast document conversion through the Cursor MCP API
- Supports markdown, HTML, PDF, DOCX, RST, LaTeX, EPUB, TXT
- Automatic path normalization for Windows compatibility
- Multiple conversion modes: string-to-string, string-to-file, file-to-file
- Automatic copyright addition to all generated documents

## Quick Installation

We provide automated installation scripts for all major operating systems:

### Windows

1. Download `install_windows.bat` from this repository
2. Right-click and select "Run as administrator"
3. Follow the on-screen instructions

### macOS

1. Download `install_macos.sh` from this repository
2. Open Terminal and navigate to the download location
3. Run: `chmod +x install_macos.sh && sudo ./install_macos.sh`
4. Follow the on-screen instructions

### Linux (Ubuntu/Debian)

1. Download `install_linux.sh` from this repository
2. Open Terminal and navigate to the download location
3. Run: `chmod +x install_linux.sh && sudo ./install_linux.sh`
4. Follow the on-screen instructions

The installation scripts will:
- Check for required dependencies
- Create the necessary directories
- Clone the repository
- Build the application
- Configure Cursor IDE

## Manual Installation

### Prerequisites

- [Pandoc](https://pandoc.org/installing.html) must be installed and available in your PATH
  - Download and install manually from the [official website](https://pandoc.org/installing.html)
  - **Important**: Restart your computer after installing Pandoc
  - You can verify installation by running `pandoc --version` in terminal
- For PDF generation, a LaTeX distribution is required (MiKTeX recommended for Windows)
- Git and Go programming language must be installed

### Manual Build

1. Clone this repository
2. Build the server:
   ```
   go build -o pandoc-server.exe ./cmd/main.go
   ```
3. Run the server:
   ```
   ./pandoc-server.exe
   ```

## Integration with Cursor IDE

To integrate with Cursor IDE, add the following to your MCP configuration file (`mcp.json`):

```json
// In the "servers" section
"pandoc_mcp_go": {
  "type": "stdio",
  "command": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go.exe",
  "cwd": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go"
}

// In the "roots" section
"pandoc": {
  "type": "document-converter",
  "server": "pandoc_mcp_go"
}
```

## Usage Examples

### Convert markdown to HTML

```go
result, err := converter.ConvertString("# Hello World", "markdown", "html")
```

### Convert markdown to Word document

```go
err := converter.ConvertStringToFile("# Hello World", "markdown", "docx", "output.docx")
```

### Convert existing file to PDF

```go
err := converter.ConvertFile("input.md", "markdown", "pdf", "output.pdf")
```

## Example Scripts

Check the `examples` directory for sample files and scripts:
- `sample.md` - Sample markdown document for testing
- `convert_example.ps1` - PowerShell script demonstrating how to use the server

## License

© 2025 SnowWhite AI - All Rights Reserved

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 