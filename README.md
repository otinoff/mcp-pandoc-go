# MCP-Pandoc

A simple and powerful document conversion server for Cursor IDE using Pandoc.

## When to use this (Jobs To Be Done)

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

## Getting Started

### Prerequisites

- [Pandoc](https://pandoc.org/installing.html) must be installed and available in your PATH
  - For Windows, you can use `install_pandoc.bat` script to install it automatically
  - Or download and install manually from the [official website](https://pandoc.org/installing.html)
- For PDF generation, a LaTeX distribution is required (MiKTeX recommended for Windows)

### Installation

#### Option 1: Using the provided build script (Windows)

1. Clone this repository
2. Run `install_pandoc.bat` if Pandoc is not installed
3. Run `build.bat` to build the server
4. The executable `pandoc-mcp-go.exe` will be created in the project directory

#### Option 2: Manual build

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