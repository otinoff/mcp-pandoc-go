# Example of using MCP-Pandoc-Go server
# This script demonstrates how to send requests for document conversion

# Check if the server is running
$processName = "pandoc-mcp-go"
$serverRunning = Get-Process -Name $processName -ErrorAction SilentlyContinue

if (-not $serverRunning) {
    Write-Host "Server $processName is not running. Starting..."
    Start-Process -FilePath "..\pandoc-mcp-go.exe" -NoNewWindow
    
    # Wait a bit for the server to start
    Start-Sleep -Seconds 2
}

# Path to sample markdown file
$sampleFile = Join-Path $PSScriptRoot "sample.md"
$outputHtml = Join-Path $PSScriptRoot "output.html"
$outputPdf = Join-Path $PSScriptRoot "output.pdf"
$outputDocx = Join-Path $PSScriptRoot "output.docx"

# Function to send request to MCP server
function Send-MCPRequest {
    param (
        [Parameter(Mandatory=$true)]
        [string]$InputFile,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputFile,
        
        [Parameter(Mandatory=$true)]
        [string]$InputFormat,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputFormat
    )
    
    $request = @{
        "method" = "tools/call"
        "params" = @{
            "name" = "convert_contents"
            "arguments" = @{
                "input_file" = $InputFile
                "output_file" = $OutputFile
                "input_format" = $InputFormat
                "output_format" = $OutputFormat
            }
        }
        "id" = 1
    } | ConvertTo-Json -Depth 5
    
    $requestObj = $request | ConvertFrom-Json
    Write-Host "Sending request: $($requestObj.method)" -ForegroundColor Yellow
    Write-Host "  Converting from $InputFormat to $OutputFormat" -ForegroundColor Cyan
    Write-Host "  Input: $InputFile" -ForegroundColor Cyan
    Write-Host "  Output: $OutputFile" -ForegroundColor Cyan
    
    try {
        # In a real scenario, we would send the request through the MCP client
        # This is just a demonstration of the request format
        Write-Host "Request in JSON format:" -ForegroundColor Green
        Write-Host $request
        
        # In a real scenario, there would be code here to send the request to the MCP server
        # For example, via stdin/stdout or HTTP, depending on the server configuration
    }
    catch {
        Write-Host "Error sending request: $_" -ForegroundColor Red
    }
}

# Convert Markdown to HTML
Write-Host "=== Converting Markdown to HTML ===" -ForegroundColor Magenta
Send-MCPRequest -InputFile $sampleFile -OutputFile $outputHtml -InputFormat "markdown" -OutputFormat "html"

# Convert Markdown to PDF (requires TeX installed)
Write-Host "`n=== Converting Markdown to PDF ===" -ForegroundColor Magenta
Send-MCPRequest -InputFile $sampleFile -OutputFile $outputPdf -InputFormat "markdown" -OutputFormat "pdf"

# Convert Markdown to DOCX
Write-Host "`n=== Converting Markdown to DOCX ===" -ForegroundColor Magenta
Send-MCPRequest -InputFile $sampleFile -OutputFile $outputDocx -InputFormat "markdown" -OutputFormat "docx"

Write-Host "`nFor use in a real scenario, connect the MCP client to this server via mcp.json"
Write-Host "See the README.md file for examples of requests" 