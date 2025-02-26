package log

import (
	defaultLog "log"
	"os"
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
