@echo off
echo Building MCP-Pandoc-Go...

REM Установка модулей
echo Installing dependencies...
go get github.com/mark3labs/mcp-go

REM Сборка для Windows
echo Building for Windows...
go build -o pandoc-mcp-go.exe ./cmd

echo Build complete: pandoc-mcp-go.exe

REM Проверка работоспособности
where pandoc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: Pandoc not found in PATH. Make sure to install Pandoc or set PANDOC_PATH environment variable.
    echo Download Pandoc from: https://pandoc.org/installing.html
)

echo.
echo To use the server, add the following to your MCP_servers/mcp.json "servers" section:
echo.
echo "pandoc_mcp_go": {
echo   "type": "stdio",
echo   "command": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go.exe",
echo   "cwd": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go"
echo }
echo.
echo And in the "roots" section:
echo.
echo "pandoc": {
echo   "type": "document-converter",
echo   "server": "pandoc_mcp_go"
echo } 