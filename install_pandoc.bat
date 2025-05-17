@echo off
echo Checking if Pandoc is installed...

where pandoc >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo Pandoc already installed.
    pandoc --version
    goto :CONFIG
) else (
    echo Pandoc not found. Downloading installer...
)

set PANDOC_VERSION=3.1.11
set INSTALLER_URL=https://github.com/jgm/pandoc/releases/download/%PANDOC_VERSION%/pandoc-%PANDOC_VERSION%-windows-x86_64.msi
set INSTALLER_PATH=%TEMP%\pandoc-%PANDOC_VERSION%-installer.msi

echo Downloading Pandoc %PANDOC_VERSION%...
powershell -Command "Invoke-WebRequest -Uri '%INSTALLER_URL%' -OutFile '%INSTALLER_PATH%'"

if %ERRORLEVEL% NEQ 0 (
    echo Failed to download Pandoc installer.
    echo Please download and install Pandoc manually from https://pandoc.org/installing.html
    goto :eof
)

echo Installing Pandoc...
msiexec /i "%INSTALLER_PATH%" /qb

if %ERRORLEVEL% NEQ 0 (
    echo Failed to install Pandoc.
    echo Please install Pandoc manually from https://pandoc.org/installing.html
    goto :eof
)

echo Pandoc installation completed. Cleaning up...
del "%INSTALLER_PATH%"

echo Verifying installation...
where pandoc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Pandoc installation failed.
    echo Please restart your computer and try again or install manually.
    goto :eof
)

pandoc --version

:CONFIG
echo.
echo Configuring mcp.json...

echo.
echo Setup complete. You can now use mcp-pandoc-go server with the Cursor IDE.
echo To build the server from source, run build.bat
echo To run the server directly, execute pandoc-mcp-go.exe

exit /b 0 