package main

import (
	"fmt"
	"log"
	"os"

	internalDB "github.com/Dpbm/jobsServer/db"
	"github.com/Dpbm/jobsServer/queue"
	serverDefinition "github.com/Dpbm/jobsServer/server"
	externalDB "github.com/Dpbm/shared/db"
	"github.com/Dpbm/shared/format"
	logger "github.com/Dpbm/shared/log"
	_ "github.com/lib/pq"
)

func main() {

	//-----LOGS------------------------------------------------------------------------------------------------------
	logFilePath := os.Getenv("LOG_FILE_PATH")
	var logFile *logger.LogFile = nil
	if logFilePath != "" {
		logFile = &logger.LogFile{}
		logFile.CreateLogFile(logFilePath) // it must execute os.Exit(1) if an error occours
		log.SetOutput(logFile.File)
	}

	//-----RABBITMQ---------------------------------------------------------------------------------------------------

	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := format.PortEnvToInt(os.Getenv("RABBITMQ_PORT")) // check port and exits if it's a number
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")

	rabbitmq := &queue.RabbitMQ{}
	rabbitmqConnection := rabbitmq.ConnectQueue(rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort) // it will exit with status 1 if an error occour
	defer rabbitmqConnection.Close()

	rabbitmqChannel := rabbitmqConnection.CreateChannel() // it will exit with status 1 if an error occour
	defer rabbitmqChannel.Close()

	//------DB---------------------------------------------------------------------------------------------------

	dbHost := os.Getenv("DB_HOST")
	dbPort := format.PortEnvToInt(os.Getenv("DB_PORT")) // in case it's not defined as an integer, the program exits
	dbUsername := os.Getenv("DB_USERNAME")
	dbPasword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbInstance := &internalDB.DB{}
	dbInstance.Connect(&externalDB.Postgres{}, dbUsername, dbPasword, dbHost, dbPort, dbName)
	defer dbInstance.CloseConnection()

	//------GRPC---------------------------------------------------------------------------------------------------

	qasmPath := os.Getenv("QASM_PATH")
	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	serverHost := os.Getenv("HOST")
	serverPort := format.PortEnvToInt(os.Getenv("PORT")) // portEnvToInt ensures that the env port is a number,
	// in other case it runs os.Exit(1)

	jobServerDefinition := &serverDefinition.JobsServer{
		QueueChannel: rabbitmqChannel,
		Database:     dbInstance,
		QasmPath:     qasmPath,
		QueueName:    queueName,
	}

	server := &serverDefinition.GRPC{}
	server.Create(serverHost, serverPort, jobServerDefinition)
	defer server.Close()

	logger.LogAction(fmt.Sprintf("Listening on host: %s", server.TCPServer.ServerURL))
	server.Listen()

	if logFile != nil {
		logFile.CloseLogFile()
	}
}
