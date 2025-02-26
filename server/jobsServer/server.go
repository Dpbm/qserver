package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	internalDB "github.com/Dpbm/jobsServer/db"
	jobsServerProto "github.com/Dpbm/jobsServer/proto"
	"github.com/Dpbm/jobsServer/queue"
	"github.com/Dpbm/jobsServer/server"
	"github.com/Dpbm/jobsServer/types"
	externalDB "github.com/Dpbm/shared/db"
	"github.com/Dpbm/shared/format"
	logger "github.com/Dpbm/shared/log"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	grpc "google.golang.org/grpc"
)

type jobsServer struct {
	jobsServerProto.UnimplementedJobsServer
	rabbitmqChannel *amqp.Channel
	database        *sql.DB
}

func checkData(data *jobsServerProto.JobProperties) error {
	if len(data.TargetSimulator) <= 0 {
		return errors.New("you must the target simulator to run your job")
	}

	// refer to: https://stackoverflow.com/questions/43770273/json-unmarshalling-without-struct
	// and: https://stackoverflow.com/questions/42152750/golang-is-there-an-easy-way-to-unmarshal-arbitrary-complex-json
	var metadata map[string]interface{}
	if data.Metadata != nil {
		err := json.Unmarshal([]byte(*data.Metadata), &metadata)

		if len(*data.Metadata) <= 0 || err != nil {
			return errors.New("invalid metadata")
		}
	}

	return nil
}

func addJobToDB(db *sql.DB, job *jobsServerProto.JobProperties, qasmFilePath string, id string) error {
	_, err := db.Exec(`
	INSERT INTO jobs(id, qasm, submission_date, target_simulator, metadata)
	VALUES ($1, $2, $3, $4, $5)
	`, id, qasmFilePath, time.Now(), job.TargetSimulator, job.Metadata)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
	INSERT INTO result_types(job_id, counts, quasi_dist, expval)
	VALUES ($1, $2, $3, $4)
	`, id, job.ResultTypeCounts, job.ResultTypeQuasiDist, job.ResultTypeExpVal)

	return err
}

func addToQASMFile(file *os.File, qasmDataChunk string) error {

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	qasmWritting, err := writer.WriteString(qasmDataChunk)

	if err != nil || qasmWritting <= 0 {
		return err
	}

	return nil
}

func addToQueue(rabbitmqChannel *amqp.Channel, jobId string) error {
	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	_, err := rabbitmqChannel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	timeoutAfter := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutAfter)
	defer cancel()

	err = rabbitmqChannel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(jobId),
		})

	if err != nil {
		return err
	}

	return nil
}

func removeFile(filename string) {
	err := os.Remove(filename)

	if err != nil {
		log.Printf("Failed on delete qasm file: %s\n", filename)
	}
}

func (server *jobsServer) AddJob(request jobsServerProto.Jobs_AddJobServer) error {
	jobData, err := request.Recv()
	if err != nil {
		return err
	}

	jobProperties := jobData.GetProperties()
	if jobProperties == nil {
		return errors.New("invalid joob properties")
	}

	err = checkData(jobProperties)
	if err != nil {
		return err
	}
	jobId := uuid.New().String()

	log.Printf("Adding new job %s\n", jobData)
	log.Printf("Job id: %s\n", jobId)

	path := os.Getenv("QASM_PATH")
	filename := jobId + ".qasm"
	qasmFilePath := filepath.Join(path, filename)

	file, err := os.Create(qasmFilePath)
	if err != nil {
		log.Fatalf("Failed on create qasm file: %s", err)
		return err
	}
	defer file.Close()

	qasmSize := 0
	for {
		req, err := request.Recv()

		if err == io.EOF {
			break
		}

		qasmChunck := req.GetQasmChunk()
		if err != nil || len(qasmChunck) <= 0 {
			log.Fatalf("Failed on get qasm data from request: %s", err)
			return err
		}

		qasmSize += len(qasmChunck)

		err = addToQASMFile(file, qasmChunck)
		if err != nil {
			log.Fatalf("Failed on add qasm data to file: %s", err)
			return err
		}
	}

	if qasmSize <= 0 {
		removeFile(qasmFilePath)
		return errors.New("you must provide qasm data")
	}

	err = addJobToDB(server.database, jobProperties, qasmFilePath, jobId)
	if err != nil {
		removeFile(qasmFilePath)
		return err
	}
	log.Println("-> Added to db")

	err = addToQueue(server.rabbitmqChannel, jobId)
	if err != nil {
		removeFile(qasmFilePath)
		return err
	}
	log.Println("-> Added to queue")

	return request.SendAndClose(&jobsServerProto.PendingJob{Id: jobId})
}

func main() {

	//--------SERVER-----------------------------------
	serverHost := os.Getenv("HOST")
	serverPort := format.PortEnvToInt(os.Getenv("PORT")) // portEnvToInt ensures that the env port is a number,
	// in other case it runs os.Exit(1)

	serverInstance := &server.Server{}
	serverInstance.Listen(serverHost, serverPort) // it will exit with status 1 if an error occour
	defer serverInstance.Close()

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

	grpcServer := grpc.NewServer()
	server := &types.JobsServer{
		QueueChannel: &rabbitmqChannel,
		Database:     dbInstance,
	}

	jobsServerProto.RegisterJobsServer(grpcServer, server)
	logger.LogAction(fmt.Sprintf("Listening on host: %s", serverInstance.ServerURL))
	grpcServer.Serve(serverInstance.Listener)
}
