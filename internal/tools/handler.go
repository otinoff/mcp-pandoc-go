package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/snowwhiteai/mcp-pandoc-go/internal/logging"
	"github.com/snowwhiteai/mcp-pandoc-go/internal/pandoc"
)

// ConvertContentsHandler handles document conversion requests
func ConvertContentsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger := logging.GetGlobalLogger()
	logger.DetailedInfo("Начало обработки запроса convert_contents")

	// Create converter
	converter, err := pandoc.NewConverter()
	if err != nil {
		logger.Error("Не удалось инициализировать Pandoc: %v", err)
		return nil, fmt.Errorf("Failed to initialize Pandoc: %v", err)
	}
	logger.Trace("Конвертер Pandoc успешно инициализирован")

	// Extract parameters
	args := req.Params.Arguments
	var contents, inputFile, inputFormat, outputFormat, outputFile string

	if val, ok := args["contents"]; ok {
		contents, _ = val.(string)
		logger.Trace("Получены входные данные в виде строки длиной %d символов", len(contents))
	}
	if val, ok := args["input_file"]; ok {
		inputFile, _ = val.(string)
		// Normalize input file path
		prevPath := inputFile
		inputFile = pandoc.NormalizePath(inputFile)
		logger.FileOperation("NORMALIZE", inputFile, true, fmt.Sprintf("Было: %s", prevPath))
	}
	if val, ok := args["input_format"]; ok {
		inputFormat, _ = val.(string)
	} else {
		inputFormat = "markdown"
	}
	if val, ok := args["output_format"]; ok {
		outputFormat, _ = val.(string)
	} else {
		outputFormat = "markdown"
	}
	if val, ok := args["output_file"]; ok {
		outputFile, _ = val.(string)
		// Normalize output file path
		prevPath := outputFile
		outputFile = pandoc.NormalizePath(outputFile)
		logger.FileOperation("NORMALIZE", outputFile, true, fmt.Sprintf("Было: %s", prevPath))
	}

	logger.DetailedInfo("Параметры конвертации: input_format=%s, output_format=%s", inputFormat, outputFormat)
	if inputFile != "" {
		logger.FileOperation("READ_INPUT", inputFile, true, "")
	}
	if outputFile != "" {
		logger.FileOperation("PREPARE_OUTPUT", outputFile, true, "")
	}

	// Check required parameters
	if contents == "" && inputFile == "" {
		logger.Error("Не указаны входные данные (contents или input_file)")
		return nil, fmt.Errorf("Either contents or input_file must be provided")
	}

	// Check that formats are valid
	if !converter.ValidateFormat(inputFormat) || !converter.ValidateFormat(outputFormat) {
		logger.Error("Неподдерживаемый формат: input=%s, output=%s", inputFormat, outputFormat)
		return nil, fmt.Errorf("Unsupported format(s): input=%s, output=%s", inputFormat, outputFormat)
	}

	// Check if PDF is used as input format (not supported by Pandoc)
	if inputFormat == "pdf" {
		logger.Error("PDF не поддерживается как входной формат для Pandoc")
		return nil, fmt.Errorf("PDF is not supported as input format, Pandoc can convert to PDF but not from PDF")
	}

	// Check if output file is needed
	needsOutputFile := outputFormat == "pdf" || outputFormat == "docx" ||
		outputFormat == "epub" || outputFormat == "latex" || outputFormat == "rst"

	if needsOutputFile && outputFile == "" {
		logger.Error("Не указан выходной файл для формата %s", outputFormat)
		return nil, fmt.Errorf("Output file is required for %s format", outputFormat)
	}

	// Check if input file exists
	if inputFile != "" {
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			logger.FileOperation("CHECK", inputFile, false, "Файл не существует")
			return nil, fmt.Errorf("Input file not found: %s", inputFile)
		}
		logger.FileOperation("CHECK", inputFile, true, "Файл существует")
	}

	// Create output directory if it doesn't exist
	if outputFile != "" {
		dir := filepath.Dir(outputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.FileOperation("CREATE_DIR", dir, false, fmt.Sprintf("Ошибка: %v", err))
			return nil, fmt.Errorf("Failed to create output directory: %v", err)
		}
		logger.FileOperation("CREATE_DIR", dir, true, "Директория создана или уже существует")
	}

	var result string
	var convertErr error

	// Run conversion based on input parameters
	if contents != "" {
		if outputFile != "" {
			// Convert string to file
			logger.Trace("Начинаем конвертацию строки в файл: %s → %s", inputFormat, outputFormat)
			convertErr = converter.ConvertStringToFile(contents, inputFormat, outputFormat, outputFile)
			if convertErr == nil {
				logger.ConversionOperation(inputFormat, outputFormat, fmt.Sprintf("Строка → %s", outputFile), true)
				result = fmt.Sprintf("Successfully converted %s to %s file: %s", inputFormat, outputFormat, outputFile)
			} else {
				logger.ConversionOperation(inputFormat, outputFormat, fmt.Sprintf("Ошибка: %v", convertErr), false)
			}
		} else {
			// Convert string to string
			logger.Trace("Начинаем конвертацию строки в строку: %s → %s", inputFormat, outputFormat)
			result, convertErr = converter.ConvertString(contents, inputFormat, outputFormat)
			if convertErr == nil {
				logger.ConversionOperation(inputFormat, outputFormat, "Строка → Строка", true)
			} else {
				logger.ConversionOperation(inputFormat, outputFormat, fmt.Sprintf("Ошибка: %v", convertErr), false)
			}
		}
	} else if inputFile != "" {
		// Convert file
		logger.Trace("Начинаем конвертацию файла: %s (%s) → %s", inputFile, inputFormat, outputFormat)
		convertErr = converter.ConvertFile(inputFile, inputFormat, outputFormat, outputFile)
		if convertErr == nil {
			if outputFile != "" {
				logger.ConversionOperation(inputFormat, outputFormat, fmt.Sprintf("%s → %s", inputFile, outputFile), true)
				result = fmt.Sprintf("Successfully converted %s to %s file: %s", inputFile, outputFormat, outputFile)
			} else {
				// Read temporary file if it was created
				tempFiles, _ := os.ReadDir(os.TempDir())
				for _, file := range tempFiles {
					if file.Name() == fmt.Sprintf("pandoc-output-*.%s", outputFormat) {
						content, err := os.ReadFile(fmt.Sprintf("%s/%s", os.TempDir(), file.Name()))
						if err == nil {
							result = string(content)
							logger.ConversionOperation(inputFormat, outputFormat, fmt.Sprintf("%s → временный файл", inputFile), true)
						} else {
							logger.FileOperation("READ", fmt.Sprintf("%s/%s", os.TempDir(), file.Name()), false, fmt.Sprintf("Ошибка: %v", err))
						}
						break
					}
				}
			}
		} else {
			logger.ConversionOperation(inputFormat, outputFormat, fmt.Sprintf("Ошибка: %v", convertErr), false)
		}
	}

	if convertErr != nil {
		logger.Error("Ошибка конвертации: %v", convertErr)
		return nil, fmt.Errorf("Conversion failed: %v", convertErr)
	}

	logger.DetailedInfo("Конвертация успешно завершена")

	// Return result depending on output type
	if needsOutputFile {
		// For binary formats return path to file
		data := map[string]string{
			"output_file": outputFile,
			"message":     result,
		}
		jsonData, _ := json.Marshal(data)
		logger.Trace("Возвращаем результат конвертации (путь к файлу): %s", outputFile)
		return mcp.NewToolResultText(string(jsonData)), nil
	} else {
		// For text formats return content
		logger.Trace("Возвращаем результат конвертации (текстовое содержимое)")
		return mcp.NewToolResultText(result), nil
	}
}
