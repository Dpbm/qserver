package log

import (
	"errors"
	"io/fs"
	defaultLog "log"
	"os"
	"path/filepath"
)

func LogError(err error) {
	defaultLog.Printf("Error: %v\n", err)
}

func LogAction(action string) {
	defaultLog.Printf("Running: %s\n", action)
}

func LogFatal(err error) {
	defaultLog.Fatalf("Error: %v\n", err)
	os.Exit(1) // ensure the program will exit with an error status
}

type LogFile struct {
	File *os.File
}

func (file *LogFile) CreateLogFile(path string) {

	_, err := os.Stat(path)

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		defaultLog.Fatalf("Error (Check path): %v\n", err)
		os.Exit(1) // ensure the program will exit with an error status
	}

	err = os.MkdirAll(filepath.Dir(path), 0770)
	if err != nil {
		defaultLog.Fatalf("Error (create path): %v\n", err)
		os.Exit(1) // ensure the program will exit with an error status
	}

	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		defaultLog.Fatalf("Error: %v\n", err)
		os.Exit(1) // ensure the program will exit with an error status
	}

	file.File = logFile
}

func (file *LogFile) CloseLogFile() {
	file.File.Close()
}
