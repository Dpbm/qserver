package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Dpbm/quantumRestAPI/db"
	dbDefinition "github.com/Dpbm/shared/db"

	"github.com/Dpbm/quantumRestAPI/server"
	"github.com/Dpbm/shared/format"
	logger "github.com/Dpbm/shared/log"
)

func main() {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	var logFile *logger.LogFile = nil
	if logFilePath != "" {
		logFile = &logger.LogFile{}
		logFile.CreateLogFile(logFilePath) // it must execute os.Exit(1) if an error occours
		log.SetOutput(logFile.File)
	}

	port := format.PortEnvToInt(os.Getenv("PORT")) // it must execute os.Exit(1) if the port is invalid

	dbHost := os.Getenv("DB_HOST")
	dbPort := format.PortEnvToInt(os.Getenv("DB_PORT")) // should execute os.Exit(1) after logging
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbInstance := db.DB{}
	dbInstance.Connect(&dbDefinition.Postgres{}, dbHost, dbPort, dbUsername, dbPassword, dbName) // on error it should exit the program with return code 1
	defer dbInstance.CloseConnection()

	server := server.SetupServer(&dbInstance)

	portString := fmt.Sprintf(":%d", port)
	err := server.Run(portString)

	if err != nil {
		logger.LogError(err)
		logFile.CloseLogFile()
		os.Exit(1)
	}

	if logFile != nil {
		logFile.CloseLogFile()
	}
}
