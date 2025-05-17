# Пример использования MCP-Pandoc-Go сервера
# Этот скрипт демонстрирует, как отправлять запросы для конвертации документов

# Проверяем, запущен ли сервер
$processName = "pandoc-mcp-go"
$serverRunning = Get-Process -Name $processName -ErrorAction SilentlyContinue

if (-not $serverRunning) {
    Write-Host "Сервер $processName не запущен. Запускаем..."
    Start-Process -FilePath "..\pandoc-mcp-go.exe" -NoNewWindow
    
    # Ждем немного, чтобы сервер успел запуститься
    Start-Sleep -Seconds 2
}

# Путь к примеру markdown-файла
$sampleFile = Join-Path $PSScriptRoot "sample.md"
$outputHtml = Join-Path $PSScriptRoot "output.html"
$outputPdf = Join-Path $PSScriptRoot "output.pdf"
$outputDocx = Join-Path $PSScriptRoot "output.docx"

# Функция для отправки запроса в MCP-сервер
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
    Write-Host "Отправка запроса: $($requestObj.method)" -ForegroundColor Yellow
    Write-Host "  Конвертация из $InputFormat в $OutputFormat" -ForegroundColor Cyan
    Write-Host "  Вход: $InputFile" -ForegroundColor Cyan
    Write-Host "  Выход: $OutputFile" -ForegroundColor Cyan
    
    try {
        # В реальном случае здесь мы бы отправили запрос через MCP-клиент
        # Это просто демонстрация формата запроса
        Write-Host "Запрос в формате JSON:" -ForegroundColor Green
        Write-Host $request
        
        # В реальном сценарии здесь был бы код для отправки запроса к MCP-серверу
        # Например, через stdin/stdout или HTTP, в зависимости от настройки сервера
    }
    catch {
        Write-Host "Ошибка отправки запроса: $_" -ForegroundColor Red
    }
}

# Конвертация Markdown в HTML
Write-Host "=== Конвертация Markdown в HTML ===" -ForegroundColor Magenta
Send-MCPRequest -InputFile $sampleFile -OutputFile $outputHtml -InputFormat "markdown" -OutputFormat "html"

# Конвертация Markdown в PDF (требует установленного TeX)
Write-Host "`n=== Конвертация Markdown в PDF ===" -ForegroundColor Magenta
Send-MCPRequest -InputFile $sampleFile -OutputFile $outputPdf -InputFormat "markdown" -OutputFormat "pdf"

# Конвертация Markdown в DOCX
Write-Host "`n=== Конвертация Markdown в DOCX ===" -ForegroundColor Magenta
Send-MCPRequest -InputFile $sampleFile -OutputFile $outputDocx -InputFormat "markdown" -OutputFormat "docx"

Write-Host "`nДля использования в реальном сценарии, подключите клиент MCP к этому серверу через mcp.json"
Write-Host "Смотрите docs\результат.md для примеров запросов" 