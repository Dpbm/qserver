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
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	jobsServerProto "github.com/Dpbm/jobsServer/proto"
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

func checkData(data *jobsServerProto.JobData) error {
	if data.NQubits <= 0 {
		return errors.New("nQubits must be greater than 0")
	}

	if len(data.Framework) <= 0 {
		return errors.New("you must provide the name of your framework/tool")
	}

	if len(data.Qasm) <= 0 {
		return errors.New("you must provide your code in qasm format")
	}

	if data.Depth <= 0 {
		return errors.New("invalid detph")
	}

	if len(data.TargetSimulator) <= 0 {
		return errors.New("you must the target simulator to run your job")
	}

	// refer to: https://stackoverflow.com/questions/43770273/json-unmarshalling-without-struct
	// and: https://stackoverflow.com/questions/42152750/golang-is-there-an-easy-way-to-unmarshal-arbitrary-complex-json
	var metadata map[string]interface{}
	if data.Metadata != nil {
		error := json.Unmarshal([]byte(*data.Metadata), &metadata)

		if len(*data.Metadata) <= 0 || error != nil {
			return errors.New("invalid metadata")
		}
	}

	return nil
}

func addJobToDB(db *sql.DB, job *jobsServerProto.JobData, qasmFilePath string, id string) error {
	_, error := db.Exec(`
	INSERT INTO jobs(id, n_qubits, framework, qasm, depth, submission_date, target_simulator, metadata)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, id, job.NQubits, job.Framework, qasmFilePath, job.Depth, time.Now(), job.TargetSimulator, job.Metadata)

	if error != nil {
		return error
	}

	_, error = db.Exec(`
	INSERT INTO result_types(job_id, counts, quasi_dist, expval)
	VALUES ($1, $2, $3, $4)
	`, id, job.ResultTypeCounts, job.ResultTypeQuasiDist, job.ResultTypeExpVal)

	return error
}

func createQASMFile(qasmData string, jobId string) (string, error) {
	path := os.Getenv("JOBS_SERVER_QASM_PATH")

	filename := jobId + ".qasm"
	qasmFilePath := filepath.Join(path, filename)

	file, error := os.Create(qasmFilePath)

	if error != nil {
		return "", error
	}

	writer := bufio.NewWriter(file)
	qasmWritting, error := writer.WriteString(qasmData)

	if error != nil || qasmWritting <= 0 {
		return "", error
	}

	writer.Flush()
	file.Close()

	return qasmFilePath, nil
}

func addToQueue(rabbitmqChannel *amqp.Channel, jobId string) error {
	queueName := os.Getenv("QUEUE_NAME")
	_, error := rabbitmqChannel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if error != nil {
		return error
	}

	timeoutAfter := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutAfter)
	defer cancel()

	error = rabbitmqChannel.PublishWithContext(
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

	if error != nil {
		return error
	}

	return nil
}

func (server *jobsServer) AddJob(request jobsServerProto.Jobs_AddJobServer) error {
	var jobStreamingData []*jobsServerProto.JobData

	for {
		data, error := request.Recv()

		if error == io.EOF {
			// when streaming reach the end of the data, it send EOF
			break
		}

		if error != nil {
			return error
		}

		jobStreamingData = append(jobStreamingData, data)
	}

	if len(jobStreamingData) <= 0 {
		return errors.New("invalid data")
	}

	jobData := jobStreamingData[0]
	error := checkData(jobData)
	if error != nil {
		return error
	}

	qasmData := jobData.Qasm
	jobId := uuid.New().String()

	log.Printf("Handling job %s\n", jobId)

	qasmFilePath, error := createQASMFile(qasmData, jobId)
	if error != nil {
		return error
	}
	log.Println("Created qasm file")

	error = addJobToDB(server.database, jobData, qasmFilePath, jobId)
	if error != nil {
		return error
	}
	log.Println("Added to db")

	error = addToQueue(server.rabbitmqChannel, jobId)
	if error != nil {
		return error
	}
	log.Println("Added to queue")

	return request.SendAndClose(&jobsServerProto.PendingJob{Id: jobId})
}

func main() {
	serverHost := os.Getenv("JOBS_SERVER_HOST")
	serverPort := os.Getenv("JOBS_SERVER_PORT")
	_, error := strconv.Atoi(serverPort)
	if error != nil {
		log.Fatalf("Invalid Port Value error : %v", error)
		panic("Invalid Port Value!")
	}

	serverUrl := fmt.Sprintf("%s:%s", serverHost, serverPort)
	listen, error := net.Listen("tcp", serverUrl)
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
		panic("Failed on listen!")
	}

	rabbitmqHost := os.Getenv("JOBS_SERVER_RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("JOBS_SERVER_RABBITMQ_PORT")
	_, error = strconv.Atoi(rabbitmqPort)
	if error != nil {
		log.Fatalf("Invalid Port Value For RabbitMQ : %v", error)
		panic("Invalid Port Value For RabbitMQ!")
	}
	rabbitmqServerUrl := fmt.Sprintf("amqp://guest:guest@%s:%s", rabbitmqHost, rabbitmqPort)
	rabbitmqConnection, error := amqp.Dial(rabbitmqServerUrl)
	if error != nil {
		log.Fatalf("Failed on connect to rabbitmq : %v", error)
		panic("Failed rabbitmq connect")
	}
	defer rabbitmqConnection.Close()

	rabbitmqChannel, error := rabbitmqConnection.Channel()
	if error != nil {
		log.Fatalf("Failed on connect to rabbitmq channel : %v", error)
		panic("Failed rabbitmq channel")
	}
	defer rabbitmqChannel.Close()

	postgresHost := os.Getenv("JOBS_SERVER_POSTGRES_HOST")
	postgresPort := os.Getenv("JOBS_SERVER_POSTGRES_PORT")
	postgresUsername := os.Getenv("JOBS_SERVER_POSTGRES_USERNAME")
	postgresPassword := os.Getenv("JOBS_SERVER_POSTGRES_PASSWORD")
	dbname := "quantum"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUsername, postgresPassword, postgresHost, postgresPort, dbname)
	db, error := sql.Open("postgres", connStr)
	if error != nil {
		log.Fatalf("Failed on connect to postgres : %v", error)
		panic("Failed connect postgres!")
	}

	grpcServer := grpc.NewServer()
	server := &jobsServer{
		rabbitmqChannel: rabbitmqChannel,
		database:        db,
	}
	jobsServerProto.RegisterJobsServer(grpcServer, server)

	log.Printf("Listening on host: %s", serverUrl)
	grpcServer.Serve(listen)
}
