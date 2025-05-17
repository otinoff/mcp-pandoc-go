@echo off
echo MCP-Pandoc-Go Installer for Windows
echo ==================================

REM Check for administrator privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo ERROR: This script requires administrator privileges.
    echo Please right-click and select "Run as administrator"
    pause
    exit /b 1
)

REM Create installation directory
echo Creating installation directory...
mkdir "C:\SnowWhiteAI" 2>nul
mkdir "C:\SnowWhiteAI\MCP_servers" 2>nul
mkdir "C:\SnowWhiteAI\MCP_servers\mcp-pandoc-go" 2>nul

REM Check for Git
where git >nul 2>&1
if %errorLevel% neq 0 (
    echo ERROR: Git is not installed. Please install Git from https://git-scm.com/download/win
    pause
    exit /b 1
)

REM Check for Go
where go >nul 2>&1
if %errorLevel% neq 0 (
    echo ERROR: Go is not installed. Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Check for Pandoc
where pandoc >nul 2>&1
if %errorLevel% neq 0 (
    echo WARNING: Pandoc is not installed or not in PATH.
    echo Please install Pandoc from https://pandoc.org/installing.html and restart your computer.
    echo After installation, run this script again.
    pause
    exit /b 1
)

REM Clone the repository
echo Cloning repository...
cd "C:\SnowWhiteAI\MCP_servers"
if exist "mcp-pandoc-go\.git" (
    echo Repository already exists, updating...
    cd mcp-pandoc-go
    git pull
) else (
    echo Cloning fresh repository...
    git clone https://github.com/otinoff/mcp-pandoc-go.git
    cd mcp-pandoc-go
)

REM Build the application
echo Building MCP-Pandoc-Go...
go build -o pandoc-mcp-go.exe ./cmd/main.go

REM Create cursor directory if it doesn't exist
echo Setting up Cursor configuration...
mkdir "%USERPROFILE%\.cursor\mcp" 2>nul
mkdir "%USERPROFILE%\.cursor\mcp\user" 2>nul

REM Create or update MCP configuration
echo {> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo   "servers": {>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo     "pandoc_mcp_go": {>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo       "type": "stdio",>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo       "command": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go/pandoc-mcp-go.exe",>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo       "cwd": "C:/SnowWhiteAI/MCP_servers/mcp-pandoc-go">> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo     }>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo   },>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo   "roots": {>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo     "pandoc": {>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo       "type": "document-converter",>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo       "server": "pandoc_mcp_go">> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo     }>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo   }>> "%USERPROFILE%\.cursor\mcp\user\config.json"
echo }>> "%USERPROFILE%\.cursor\mcp\user\config.json"

echo.
echo Installation completed successfully!
echo MCP-Pandoc-Go has been installed to: C:\SnowWhiteAI\MCP_servers\mcp-pandoc-go
echo Configuration has been set up for Cursor IDE.
echo.
echo Please restart Cursor if it's already running.
echo.
pause 