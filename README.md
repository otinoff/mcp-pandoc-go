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
- Or set `PANDOC_PATH` environment variable to the Pandoc executable

### Installation

1. Clone this repository
2. Build the server:
   ```
   go build -o pandoc-server.exe ./cmd/main.go
   ```
3. Run the server:
   ```
   ./pandoc-server.exe
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

## License

© 2025 SnowWhite AI - All Rights Reserved

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 