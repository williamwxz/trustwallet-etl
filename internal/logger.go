package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
)

func InitLogger(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	Info = log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(f, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
