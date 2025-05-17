package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/snowwhiteai/mcp-pandoc-go/internal/pandoc"
)

// ConvertContentsHandler handles document conversion requests
func ConvertContentsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create converter
	converter, err := pandoc.NewConverter()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Pandoc: %v", err)
	}

	// Extract parameters
	args := req.Params.Arguments
	var contents, inputFile, inputFormat, outputFormat, outputFile string

	if val, ok := args["contents"]; ok {
		contents, _ = val.(string)
	}
	if val, ok := args["input_file"]; ok {
		inputFile, _ = val.(string)
		// Normalize input file path
		inputFile = pandoc.NormalizePath(inputFile)
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
		outputFile = pandoc.NormalizePath(outputFile)
	}

	// Check required parameters
	if contents == "" && inputFile == "" {
		return nil, fmt.Errorf("Either contents or input_file must be provided")
	}

	// Check that formats are valid
	if !converter.ValidateFormat(inputFormat) || !converter.ValidateFormat(outputFormat) {
		return nil, fmt.Errorf("Unsupported format(s): input=%s, output=%s", inputFormat, outputFormat)
	}

	// Check if output file is needed
	needsOutputFile := outputFormat == "pdf" || outputFormat == "docx" ||
		outputFormat == "epub" || outputFormat == "latex" || outputFormat == "rst"

	if needsOutputFile && outputFile == "" {
		return nil, fmt.Errorf("Output file is required for %s format", outputFormat)
	}

	var result string
	var convertErr error

	// Run conversion based on input parameters
	if contents != "" {
		if outputFile != "" {
			// Convert string to file
			convertErr = converter.ConvertStringToFile(contents, inputFormat, outputFormat, outputFile)
			if convertErr == nil {
				result = fmt.Sprintf("Successfully converted %s to %s file: %s", inputFormat, outputFormat, outputFile)
			}
		} else {
			// Convert string to string
			result, convertErr = converter.ConvertString(contents, inputFormat, outputFormat)
		}
	} else if inputFile != "" {
		// Convert file
		convertErr = converter.ConvertFile(inputFile, inputFormat, outputFormat, outputFile)
		if convertErr == nil {
			if outputFile != "" {
				result = fmt.Sprintf("Successfully converted %s to %s file: %s", inputFile, outputFormat, outputFile)
			} else {
				// Read temporary file if it was created
				tempFiles, _ := os.ReadDir(os.TempDir())
				for _, file := range tempFiles {
					if file.Name() == fmt.Sprintf("pandoc-output-*.%s", outputFormat) {
						content, _ := os.ReadFile(fmt.Sprintf("%s/%s", os.TempDir(), file.Name()))
						result = string(content)
						break
					}
				}
			}
		}
	}

	if convertErr != nil {
		return nil, fmt.Errorf("Conversion failed: %v", convertErr)
	}

	// Return result depending on output type
	if needsOutputFile {
		// For binary formats return path to file
		data := map[string]string{
			"output_file": outputFile,
			"message":     result,
		}
		jsonData, _ := json.Marshal(data)
		return mcp.NewToolResultText(string(jsonData)), nil
	} else {
		// For text formats return content
		return mcp.NewToolResultText(result), nil
	}
}
