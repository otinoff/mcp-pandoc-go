package main

import (
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/snowwhiteai/mcp-pandoc-go/internal/logging"
	"github.com/snowwhiteai/mcp-pandoc-go/internal/tools"
)

func main() {
	// Initialize logging
	logger := logging.NewLogger("[MCP-Pandoc] ", os.Stderr)
	logger.Info("Starting MCP-Pandoc server")

	// Create MCP server
	s := server.NewMCPServer(
		"Pandoc Document Converter",
		"1.0.0",
		server.WithLogging(),
	)

	// Register convert_contents tool
	convertTool := mcp.NewTool("convert_contents",
		mcp.WithDescription("Convert document between different formats using Pandoc"),
		mcp.WithString("contents",
			mcp.Description("Source content to convert (required if input_file not provided)"),
		),
		mcp.WithString("input_file",
			mcp.Description("Complete path to input file (required if contents not provided)"),
		),
		mcp.WithString("input_format",
			mcp.Description("Source format of the content"),
			mcp.DefaultString("markdown"),
			mcp.Enum("markdown", "html", "pdf", "docx", "rst", "latex", "epub", "txt"),
		),
		mcp.WithString("output_format",
			mcp.Description("Target format"),
			mcp.DefaultString("markdown"),
			mcp.Enum("markdown", "html", "pdf", "docx", "rst", "latex", "epub", "txt"),
		),
		mcp.WithString("output_file",
			mcp.Description("Complete path for output file (required for pdf, docx, rst, latex, epub formats)"),
		),
	)

	// Add tool handler
	s.AddTool(convertTool, tools.ConvertContentsHandler)

	// Start server via stdio
	logger.Info("Server initialized, waiting for requests...")
	if err := server.ServeStdio(s); err != nil {
		logger.Error("Server error: %v", err)
		os.Exit(1)
	}
}
