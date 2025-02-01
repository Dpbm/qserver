package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	jobsServerProto "github.com/Dpbm/jobsServer/proto"
	"github.com/google/uuid"
	grpc "google.golang.org/grpc"
)

type jobsServer struct {
	jobsServerProto.UnimplementedJobsServer
}

func checkData(data *jobsServerProto.JobData) error {
	if data.NQubits <= 0 {
		return errors.New("nQubits must be greater than 0")
	}

	if len(data.Framework) <= 0 {
		return errors.New("your must provide the name of your framework/tool")
	}

	if len(data.Qasm) <= 0 {
		return errors.New("your must provide your code in qasm format")
	}

	if data.Depth <= 0 {
		return errors.New("invalid detph")
	}

	if len(data.ResultsTypes) <= 0 {
		return errors.New("your must provide the result types you want to retrieve from your job")
	}

	if len(data.TargetSimulator) <= 0 {
		return errors.New("your must the target simulator to run your job")
	}

	if data.Metadata != nil && len(*data.Metadata) <= 0 {
		return errors.New("invalid metadata")
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
			log.Fatalf("Failed receive: %v", error)
			return error
		}

		jobStreamingData = append(jobStreamingData, data)
	}

	if len(jobStreamingData) <= 0 {
		return errors.New("invalid data")
	}

	jobData := jobStreamingData[0]
	dataError := checkData(jobData)
	if dataError != nil {
		return dataError
	}

	// test volume files sharing accross multiple containers
	// send data to postgres db
	// submit job to rabbitmq (pass only jobId)

	qasmData := jobData.Qasm

	path := os.Getenv("JOBS_SERVER_QASM_PATH")
	jobId := uuid.New().String()

	filename := jobId + ".qasm"
	qasmFilePath := filepath.Join(path, filename)

	file, error := os.Create(qasmFilePath)

	if error != nil {
		return error
	}

	writer := bufio.NewWriter(file)
	qasmWritting, error := writer.WriteString(qasmData)

	if error != nil || qasmWritting <= 0 {
		return error
	}

	writer.Flush()
	file.Close()

	return request.SendAndClose(&jobsServerProto.PendingJob{Id: jobId})
}

func main() {
	host := os.Getenv("JOBS_SERVER_HOST")
	port := os.Getenv("JOBS_SERVER_PORT")

	_, error := strconv.Atoi(port)

	if error != nil {
		log.Fatalf("Invalid Port Value error : %v", error)
		panic("Invalid Port Value!")
	}

	serverUrl := fmt.Sprintf("%s:%s", host, port)

	listen, error := net.Listen("tcp", serverUrl)

	if error != nil {
		log.Fatalf("failed to listen: %v", error)
		panic("Failed on listen!")
	}

	grpcServer := grpc.NewServer()
	jobsServerProto.RegisterJobsServer(grpcServer, &jobsServer{})

	log.Printf("Listening on host: %s", host)
	grpcServer.Serve(listen)
}
