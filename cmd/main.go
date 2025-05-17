package main

import (
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/snowwhiteai/mcp-pandoc-go/internal/logging"
	"github.com/snowwhiteai/mcp-pandoc-go/internal/tools"
)

func main() {
	// Проверяем и устанавливаем директорию логов
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		// Если переменная не установлена, используем директорию логов внутри проекта
		exePath, err := os.Executable()
		if err == nil {
			exeDir := filepath.Dir(exePath)
			logDir = filepath.Join(exeDir, "logs")
		} else {
			// Если не удалось получить путь к исполняемому файлу, используем текущую директорию
			logDir = "logs"
		}
		os.Setenv("LOG_DIR", logDir)
	}

	// Убедимся, что директория логов существует
	if err := os.MkdirAll(logDir, 0755); err != nil {
		os.Stderr.WriteString("ERROR: Failed to create log directory: " + err.Error() + "\n")
	}

	// Initialize logging
	logger := logging.NewLogger("[MCP-Pandoc] ", os.Stderr)
	logging.InitGlobalLogger("[MCP-Pandoc] ", os.Stderr)
	logger.Info("Starting MCP-Pandoc server")
	logger.Info("Log directory set to: %s", logDir)

	// Set detailed logging level if specified
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		logger.Info("Log level set to %s", logLevel)
	} else {
		// По умолчанию устанавливаем уровень debug для более подробного логирования
		os.Setenv("LOG_LEVEL", "debug")
		logLevel = "debug"
		logger.Info("Log level defaulted to debug")
	}

	// Запись тестового сообщения для проверки логирования
	logger.DetailedInfo("=== LOGGING TEST: Detailed Info message ===")
	logger.Debug("=== LOGGING TEST: Debug message ===")
	logger.Trace("=== LOGGING TEST: Trace message ===")
	logger.FileOperation("TEST", logDir, true, "Проверка записи операций с файлами")
	logger.ConversionOperation("test", "test", "Проверка записи операций конвертации", true)

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
