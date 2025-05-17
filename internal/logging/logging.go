package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger представляет систему логирования
type Logger struct {
	logger        *log.Logger
	file          *os.File
	isFileLogging bool
}

// NewLogger создает новый логгер с выводом в stderr и файл (если задана переменная LOG_DIR)
func NewLogger(prefix string, w io.Writer) *Logger {
	logger := &Logger{
		logger:        log.New(w, prefix, log.LstdFlags),
		isFileLogging: false,
	}

	// Проверяем, нужно ли логировать в файл
	logDir := os.Getenv("LOG_DIR")
	if logDir != "" {
		// Создаем директорию для логов, если она не существует
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
			return logger
		}

		// Создаем файл лога с датой в имени
		logFile := filepath.Join(logDir, fmt.Sprintf("pandoc-mcp-%s.log", time.Now().Format("2006-01-02")))
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
			return logger
		}

		// Создаем логгер с выводом и в stderr, и в файл
		mw := io.MultiWriter(w, f)
		logger.logger = log.New(mw, prefix, log.LstdFlags)
		logger.file = f
		logger.isFileLogging = true

		// Записываем начальное сообщение
		logger.logger.Printf("INFO: ------------ НАЧАЛО СЕССИИ ЛОГИРОВАНИЯ ------------")
		logger.logger.Printf("INFO: Logging to %s", logFile)

		// Принудительная запись на диск
		if logger.file != nil {
			logger.file.Sync()
		}
	}

	return logger
}

// записать в файл с принудительной синхронизацией
func (l *Logger) write(level, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Printf("%s: %s", level, msg)

	// Принудительно сбрасываем буфер, если логируем в файл
	if l.isFileLogging && l.file != nil {
		l.file.Sync()
	}
}

// Info логирует информационное сообщение
func (l *Logger) Info(format string, v ...interface{}) {
	l.write("INFO", format, v...)
}

// Error логирует сообщение об ошибке
func (l *Logger) Error(format string, v ...interface{}) {
	l.write("ERROR", format, v...)
}

// Debug логирует отладочное сообщение (только если LOG_LEVEL=debug)
func (l *Logger) Debug(format string, v ...interface{}) {
	if os.Getenv("LOG_LEVEL") == "debug" || os.Getenv("LOG_LEVEL") == "trace" {
		l.write("DEBUG", format, v...)
	}
}

// Trace логирует детальное сообщение о ходе выполнения операции (LOG_LEVEL=trace или debug)
func (l *Logger) Trace(format string, v ...interface{}) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "trace" || logLevel == "debug" {
		l.write("TRACE", format, v...)
	}
}

// DetailedInfo логирует подробное информационное сообщение независимо от уровня логирования
func (l *Logger) DetailedInfo(format string, v ...interface{}) {
	l.write("DETAIL", format, v...)
}

// FileOperation логирует операции с файлами (проверка, создание, чтение, запись)
func (l *Logger) FileOperation(operation, path string, success bool, details string) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}

	msg := fmt.Sprintf("FILE OP: %s [%s] Path: %s", operation, status, path)
	if details != "" {
		msg += fmt.Sprintf(" - %s", details)
	}

	l.write("OPERATION", msg)
}

// ConversionOperation логирует операции конвертации
func (l *Logger) ConversionOperation(inputFormat, outputFormat, details string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}

	msg := fmt.Sprintf("CONVERT: %s→%s [%s]", inputFormat, outputFormat, status)
	if details != "" {
		msg += fmt.Sprintf(" - %s", details)
	}

	l.write("OPERATION", msg)
}

// Close закрывает файл логов
func (l *Logger) Close() {
	if l.isFileLogging && l.file != nil {
		// Записываем завершающее сообщение
		l.logger.Printf("INFO: ------------ КОНЕЦ СЕССИИ ЛОГИРОВАНИЯ ------------")
		// Синхронизация и закрытие
		l.file.Sync()
		l.file.Close()
	}
}

// Глобальный логгер для статических функций
var globalLogger *Logger

// InitGlobalLogger инициализирует глобальный логгер
func InitGlobalLogger(prefix string, w io.Writer) {
	globalLogger = NewLogger(prefix, w)
}

// GetGlobalLogger возвращает глобальный логгер
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewLogger("[MCP-Pandoc] ", os.Stderr)
	}
	return globalLogger
}
