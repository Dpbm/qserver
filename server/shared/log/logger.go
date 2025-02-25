package log

import (
	defaultLog "log"
)

func LogError(err error) {
	defaultLog.Printf("Error: %v\n", err)
}

func LogAction(action string) {
	defaultLog.Printf("Running: %s\n", action)
}

func LogFatal(err error) {
	defaultLog.Fatalf("Error: %v\n", err)
}
