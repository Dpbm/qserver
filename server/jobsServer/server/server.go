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

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func checkData(data *jobsServerProto.JobData) error {
	if len(data.Qasm) <= 0 {
		return errors.New("you must provide your code in qasm format")
	}

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

func addJobToDB(db *sql.DB, job *jobsServerProto.JobData, qasmFilePath string, id string) error {
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

func createQASMFile(qasmData string, jobId string) (string, error) {
	path := os.Getenv("QASM_PATH")

	filename := jobId + ".qasm"
	qasmFilePath := filepath.Join(path, filename)

	file, err := os.Create(qasmFilePath)

	if err != nil {
		return "", err
	}

	writer := bufio.NewWriter(file)
	qasmWritting, err := writer.WriteString(qasmData)

	if err != nil || qasmWritting <= 0 {
		return "", err
	}

	writer.Flush()
	file.Close()

	return qasmFilePath, nil
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

	message := Message{
		Type: "job",
		Data: jobId,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = rabbitmqChannel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonData,
		})

	if err != nil {
		return err
	}

	return nil
}

func (server *jobsServer) AddJob(request jobsServerProto.Jobs_AddJobServer) error {
	var jobStreamingData []*jobsServerProto.JobData

	for {
		data, err := request.Recv()

		if err == io.EOF {
			// when streaming reach the end of the data, it send EOF
			break
		}

		if err != nil {
			return err
		}

		jobStreamingData = append(jobStreamingData, data)
	}

	if len(jobStreamingData) <= 0 {
		return errors.New("invalid data")
	}

	jobData := jobStreamingData[0]
	err := checkData(jobData)
	if err != nil {
		return err
	}

	qasmData := jobData.Qasm
	jobId := uuid.New().String()

	log.Printf("Handling job %s\n", jobId)

	qasmFilePath, err := createQASMFile(qasmData, jobId)
	if err != nil {
		return err
	}
	log.Println("Created qasm file")

	err = addJobToDB(server.database, jobData, qasmFilePath, jobId)
	if err != nil {
		return err
	}
	log.Println("Added to db")

	err = addToQueue(server.rabbitmqChannel, jobId)
	if err != nil {
		return err
	}
	log.Println("Added to queue")

	return request.SendAndClose(&jobsServerProto.PendingJob{Id: jobId})
}

func main() {
	serverHost := os.Getenv("HOST")
	serverPort := os.Getenv("PORT")
	_, err := strconv.Atoi(serverPort)
	if err != nil {
		log.Fatalf("Invalid Port Value error : %v", err)
		panic("Invalid Port Value!")
	}

	serverUrl := fmt.Sprintf("%s:%s", serverHost, serverPort)
	listen, err := net.Listen("tcp", serverUrl)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic("Failed on listen!")
	}

	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	_, err = strconv.Atoi(rabbitmqPort)
	if err != nil {
		log.Fatalf("Invalid Port Value For RabbitMQ : %v", err)
		panic("Invalid Port Value For RabbitMQ!")
	}
	rabbitmqServerUrl := fmt.Sprintf("amqp://guest:guest@%s:%s", rabbitmqHost, rabbitmqPort)
	rabbitmqConnection, err := amqp.Dial(rabbitmqServerUrl)
	if err != nil {
		log.Fatalf("Failed on connect to rabbitmq : %v", err)
		panic("Failed rabbitmq connect")
	}
	defer rabbitmqConnection.Close()

	rabbitmqChannel, err := rabbitmqConnection.Channel()
	if err != nil {
		log.Fatalf("Failed on connect to rabbitmq channel : %v", err)
		panic("Failed rabbitmq channel")
	}
	defer rabbitmqChannel.Close()

	postgresHost := os.Getenv("DB_HOST")
	postgresPort := os.Getenv("DB_PORT")
	_, err = strconv.Atoi(postgresPort)
	if err != nil {
		log.Fatalf("Invalid Port Value For DB : %v", err)
		panic("Invalid Port Value For DB!")
	}
	postgresUsername := os.Getenv("DB_USERNAME")
	postgresPassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUsername, postgresPassword, postgresHost, postgresPort, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed on connect to postgres : %v", err)
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
