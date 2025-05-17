@echo off
echo ======================================
echo      Building MCP-Pandoc-Go...
echo ======================================

REM Установка модулей
echo [1/3] Installing dependencies...
go get github.com/mark3labs/mcp-go
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to install dependencies.
    echo.
    echo Build process FAILED!
    exit /b 1
)

REM Сборка для Windows
echo [2/3] Building for Windows...
go build -o pandoc-mcp-go.exe ./cmd
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Build compilation failed.
    echo.
    echo Build process FAILED!
    exit /b 1
)

REM Проверка наличия собранного исполняемого файла
if not exist "pandoc-mcp-go.exe" (
    echo ERROR: Build executable not found.
    echo.
    echo Build process FAILED!
    exit /b 1
)

echo [3/3] Verifying Pandoc installation...
where pandoc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: Pandoc not found in PATH. Make sure to install Pandoc or set PANDOC_PATH environment variable.
    echo Download Pandoc from: https://pandoc.org/installing.html
)

echo.
echo ======================================
echo      BUILD COMPLETED SUCCESSFULLY
echo ======================================
echo.
echo Output: pandoc-mcp-go.exe
echo.

REM Вспомогательная информация по интеграции с MCP
echo ------------------- INFO -------------------
echo To configure the server in MCP, add to your mcp.json:
echo.
echo "servers" section:
echo "pandoc_mcp_go": {
echo   "type": "stdio",
echo   "command": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go.exe",
echo   "cwd": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go"
echo }
echo.
echo "roots" section:
echo "pandoc": {
echo   "type": "document-converter",
echo   "server": "pandoc_mcp_go"
echo }
echo ------------------------------------------- 