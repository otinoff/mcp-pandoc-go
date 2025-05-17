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
	logger *log.Logger
	file   *os.File
}

// NewLogger создает новый логгер с выводом в stderr и файл (если задана переменная LOG_DIR)
func NewLogger(prefix string, w io.Writer) *Logger {
	logger := &Logger{
		logger: log.New(w, prefix, log.LstdFlags),
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

		logger.Info("Logging to %s", logFile)
	}

	return logger
}

// Info логирует информационное сообщение
func (l *Logger) Info(format string, v ...interface{}) {
	l.logger.Printf("INFO: "+format, v...)
}

// Error логирует сообщение об ошибке
func (l *Logger) Error(format string, v ...interface{}) {
	l.logger.Printf("ERROR: "+format, v...)
}

// Debug логирует отладочное сообщение (только если LOG_LEVEL=debug)
func (l *Logger) Debug(format string, v ...interface{}) {
	if os.Getenv("LOG_LEVEL") == "debug" {
		l.logger.Printf("DEBUG: "+format, v...)
	}
}

// Close закрывает файл логов
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
