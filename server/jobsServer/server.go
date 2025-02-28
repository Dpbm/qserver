package main

import (
	"fmt"
	"os"

	internalDB "github.com/Dpbm/jobsServer/db"
	"github.com/Dpbm/jobsServer/queue"
	"github.com/Dpbm/jobsServer/server"
	serverDefinition "github.com/Dpbm/jobsServer/server"
	externalDB "github.com/Dpbm/shared/db"
	"github.com/Dpbm/shared/format"
	logger "github.com/Dpbm/shared/log"
	_ "github.com/lib/pq"
)

func main() {

	//-----RABBITMQ---------------------------------------------------------------------------------------------------

	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := format.PortEnvToInt(os.Getenv("RABBITMQ_PORT")) // check port and exits if it's a number

	rabbitmq := &queue.RabbitMQ{}
	rabbitmqConnection := rabbitmq.ConnectQueue(rabbitmqHost, rabbitmqPort, "guest", "guest") // it will exit with status 1 if an error occour
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

	server := &server.GRPC{}
	server.Create(serverHost, serverPort, jobServerDefinition)
	defer server.Close()

	logger.LogAction(fmt.Sprintf("Listening on host: %s", server.TCPServer.ServerURL))
	server.Listen()
}
