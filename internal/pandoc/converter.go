package pandoc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// PandocConverter represents a document converter based on Pandoc
type PandocConverter struct {
	pandocPath string
}

// NewConverter creates a new document converter
func NewConverter() (*PandocConverter, error) {
	// Try to get path from environment variable
	pandocPath := os.Getenv("PANDOC_PATH")
	if pandocPath == "" {
		// If variable is not set, look for pandoc in system PATH
		path, err := exec.LookPath("pandoc")
		if err != nil {
			return nil, errors.New("Pandoc not found in PATH, please set PANDOC_PATH environment variable")
		}
		pandocPath = path
	}

	// Check that file exists and is executable
	info, err := os.Stat(pandocPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing Pandoc at %s: %v", pandocPath, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a Pandoc executable", pandocPath)
	}

	// On Windows don't check execution permissions
	if os.Getenv("OS") != "Windows_NT" && info.Mode()&0111 == 0 {
		return nil, fmt.Errorf("%s is not executable", pandocPath)
	}

	return &PandocConverter{
		pandocPath: pandocPath,
	}, nil
}

// normalizePath converts a file path to a format suitable for the current OS
func normalizePath(path string) string {
	// For Windows: if path starts with a slash, remove it
	if runtime.GOOS == "windows" {
		// Remove initial slash if present
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}

		// If path is like "/c:/path" or "/C:/path", convert to "C:/path"
		if len(path) >= 3 && path[1] == ':' && (path[0] >= 'a' && path[0] <= 'z' || path[0] >= 'A' && path[0] <= 'Z') {
			path = path[0:1] + ":" + path[2:]
		}

		// Replace mixed slashes with standard ones for Windows
		path = filepath.FromSlash(path)
	}

	return path
}

// NormalizePath exported version of normalizePath for use from other packages
func NormalizePath(path string) string {
	return normalizePath(path)
}

// addCopyright adds SnowWhite AI copyright to the end of the document
func addCopyright(content string) string {
	// Use simple markdown syntax instead of HTML
	copyrightText := `


--------------------------------------------------------------------------------

**© 2025 SnowWhite AI - All Rights Reserved**

`
	return content + copyrightText
}

// ValidateFormat checks if the format is supported
func (p *PandocConverter) ValidateFormat(format string) bool {
	supportedFormats := []string{"markdown", "html", "pdf", "docx", "rst", "latex", "epub", "txt"}
	for _, f := range supportedFormats {
		if f == format {
			return true
		}
	}
	return false
}

// ConvertString converts a string from one format to another
func (p *PandocConverter) ConvertString(content, inputFormat, outputFormat string) (string, error) {
	// Add copyright to source content
	if inputFormat == "markdown" {
		content = addCopyright(content)
	}

	// Format validation
	if !p.ValidateFormat(inputFormat) || !p.ValidateFormat(outputFormat) {
		return "", fmt.Errorf("unsupported format: input=%s, output=%s", inputFormat, outputFormat)
	}

	// For formats requiring a file output, return error
	if outputFormat == "pdf" || outputFormat == "docx" || outputFormat == "epub" {
		return "", fmt.Errorf("output_file is required for %s format", outputFormat)
	}

	// Create temporary file for input data
	tmpInput, err := os.CreateTemp("", "pandoc-input-*."+inputFormat)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpInput.Name())

	if _, err := tmpInput.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %v", err)
	}
	tmpInput.Close()

	// Run pandoc
	cmd := exec.Command(p.pandocPath,
		"-f", inputFormat,
		"-t", outputFormat,
		tmpInput.Name())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pandoc conversion failed: %v\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// ConvertFile converts a file from one format to another
func (p *PandocConverter) ConvertFile(inputFile, inputFormat, outputFormat, outputFile string) error {
	// Normalize paths
	inputFile = normalizePath(inputFile)
	if outputFile != "" {
		outputFile = normalizePath(outputFile)
	}

	// Format validation
	if !p.ValidateFormat(inputFormat) || !p.ValidateFormat(outputFormat) {
		return fmt.Errorf("unsupported format: input=%s, output=%s", inputFormat, outputFormat)
	}

	// Check existence of input file
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", inputFile)
	}

	// For some formats, output file must be specified
	needsOutputFile := outputFormat == "pdf" || outputFormat == "docx" ||
		outputFormat == "epub" || outputFormat == "latex" || outputFormat == "rst"

	if needsOutputFile && outputFile == "" {
		return fmt.Errorf("output_file is required for %s format", outputFormat)
	}

	// Determine path to footer file relative to executable
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	execDir := filepath.Dir(execPath)
	footerPath := filepath.Join(execDir, "templates", "footer.md")

	// Check if footer.md file exists
	footerExists := false
	if _, err := os.Stat(footerPath); err == nil {
		footerExists = true
	} else {
		// Try to find file relative to current directory
		footerPath = filepath.Join("templates", "footer.md")
		if _, err := os.Stat(footerPath); err == nil {
			footerExists = true
		}
	}

	// If input format is markdown and footer.md is not applied, add copyright to content
	if inputFormat == "markdown" && (!footerExists || outputFormat == "txt") {
		// Read input file content
		content, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %v", err)
		}

		// Add copyright
		contentWithCopyright := addCopyright(string(content))

		// Create temporary file with updated content
		tmpInput, err := os.CreateTemp("", "pandoc-input-with-copyright-*."+inputFormat)
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %v", err)
		}
		defer os.Remove(tmpInput.Name())

		if _, err := tmpInput.WriteString(contentWithCopyright); err != nil {
			return fmt.Errorf("failed to write to temporary file: %v", err)
		}
		tmpInput.Close()

		// Use temporary file instead of original
		inputFile = tmpInput.Name()
	}

	// If output_file is not specified, use temporary file
	var tmpOutputFile string
	if outputFile == "" {
		tmp, err := os.CreateTemp("", "pandoc-output-*."+outputFormat)
		if err != nil {
			return fmt.Errorf("failed to create temporary output file: %v", err)
		}
		tmpOutputFile = tmp.Name()
		tmp.Close()
		outputFile = tmpOutputFile
		defer os.Remove(tmpOutputFile)
	} else {
		// Create directory for output file if it doesn't exist
		dir := filepath.Dir(outputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Run pandoc
	args := []string{
		"-f", inputFormat,
		"-t", outputFormat,
		"-o", outputFile,
	}

	// If footer file exists and output format supports inclusions, add it
	if footerExists && (outputFormat == "docx" || outputFormat == "pdf" || outputFormat == "html") {
		args = append(args, "--include-after-body", footerPath)
	}

	// Add input file at the end
	args = append(args, inputFile)

	cmd := exec.Command(p.pandocPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pandoc conversion failed: %v\nOutput: %s", err, string(output))
	}

	// If temporary file was used, read its content
	if tmpOutputFile != "" {
		content, err := ioutil.ReadFile(tmpOutputFile)
		if err != nil {
			return fmt.Errorf("failed to read output file: %v", err)
		}
		// For simple formats return content
		if strings.Contains("markdown html txt", outputFormat) {
			fmt.Println(string(content))
		}
	}

	return nil
}

// ConvertStringToFile converts a string to a file
func (p *PandocConverter) ConvertStringToFile(content, inputFormat, outputFormat, outputFile string) error {
	// Normalize output file path
	if outputFile != "" {
		outputFile = normalizePath(outputFile)
	}

	// Format validation
	if !p.ValidateFormat(inputFormat) || !p.ValidateFormat(outputFormat) {
		return fmt.Errorf("unsupported format: input=%s, output=%s", inputFormat, outputFormat)
	}

	// For some formats, output file must be specified
	needsOutputFile := outputFormat == "pdf" || outputFormat == "docx" ||
		outputFormat == "epub" || outputFormat == "latex" || outputFormat == "rst"

	if needsOutputFile && outputFile == "" {
		return fmt.Errorf("output_file is required for %s format", outputFormat)
	}

	// If input format is markdown, add copyright
	if inputFormat == "markdown" {
		content = addCopyright(content)
	}

	// Create temporary file for input data
	tmpInput, err := os.CreateTemp("", "pandoc-input-*."+inputFormat)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpInput.Name())

	if _, err := tmpInput.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to temporary file: %v", err)
	}
	tmpInput.Close()

	// Determine path to footer file relative to executable
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	execDir := filepath.Dir(execPath)
	footerPath := filepath.Join(execDir, "templates", "footer.md")

	// Check if footer.md file exists
	footerExists := false
	if _, err := os.Stat(footerPath); err == nil {
		footerExists = true
	} else {
		// Try to find file relative to current directory
		footerPath = filepath.Join("templates", "footer.md")
		if _, err := os.Stat(footerPath); err == nil {
			footerExists = true
		}
	}

	// If output_file is not specified for formats that can be returned as text
	if outputFile == "" {
		// Create temporary file for output
		tmpOutput, err := os.CreateTemp("", "pandoc-output-*."+outputFormat)
		if err != nil {
			return fmt.Errorf("failed to create temporary output file: %v", err)
		}
		outputFile = tmpOutput.Name()
		tmpOutput.Close()
		defer os.Remove(outputFile)
	} else {
		// Create directory for output file if it doesn't exist
		dir := filepath.Dir(outputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Run pandoc
	args := []string{
		"-f", inputFormat,
		"-t", outputFormat,
		"-o", outputFile,
	}

	// If footer file exists and output format supports inclusions, add it
	if footerExists && (outputFormat == "docx" || outputFormat == "pdf" || outputFormat == "html") {
		args = append(args, "--include-after-body", footerPath)
	}

	// Add input file at the end
	args = append(args, tmpInput.Name())

	cmd := exec.Command(p.pandocPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pandoc conversion failed: %v\nOutput: %s", err, string(output))
	}

	// If output file was not specified and format allows text output
	if !needsOutputFile && outputFile != "" {
		content, err := ioutil.ReadFile(outputFile)
		if err != nil {
			return fmt.Errorf("failed to read output file: %v", err)
		}
		fmt.Println(string(content))
	}

	return nil
}
