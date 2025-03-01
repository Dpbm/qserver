package server

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Dpbm/jobsServer/data"
	"github.com/Dpbm/jobsServer/db"
	"github.com/Dpbm/jobsServer/files"
	jobsServerProto "github.com/Dpbm/jobsServer/proto"
	"github.com/Dpbm/jobsServer/queue"
	"github.com/Dpbm/shared/log"
	logger "github.com/Dpbm/shared/log"
	"github.com/google/uuid"
)

type JobsServer struct {
	jobsServerProto.UnimplementedJobsServer
	QueueChannel queue.QueueChannel
	Database     *db.DB
	QasmPath     string
	QueueName    string
}

func (server *JobsServer) AddJob(request jobsServerProto.Jobs_AddJobServer) error {
	jobData, err := request.Recv()
	if err != nil {
		return err
	}

	jobProperties := jobData.GetProperties()
	err = data.CheckData(jobProperties)
	if err != nil {
		return err
	}

	jobId := uuid.New().String()
	logger.LogAction(fmt.Sprintf("Adding new job %s\n", jobData))
	logger.LogAction(fmt.Sprintf("Creating job with ID %s\n", jobId))

	filename := jobId + ".qasm"
	qasmFilePath := filepath.Join(server.QasmPath, filename)
	qasmFile := &files.Qasm{Path: qasmFilePath, Filename: filename}

	err = qasmFile.CreateFile()
	if err != nil {
		logger.LogError(err)
		return err
	}
	defer qasmFile.Close()

	log.LogAction("Getting QASM...")
	qasmSize := 0
	for {
		req, err := request.Recv()

		if err == io.EOF {
			break
		}

		qasmChunck := req.GetQasmChunk()
		if err != nil || len(qasmChunck) <= 0 {
			logger.LogError(err)
			return err
		}

		qasmSize += len(qasmChunck)

		err = qasmFile.AddChunckToFile(qasmChunck)
		if err != nil {
			logger.LogError(err)
			return err
		}
	}

	log.LogAction(fmt.Sprintf("QASM size %d", qasmSize))

	if qasmSize <= 0 {
		err = errors.New("you must provide qasm data")
		logger.LogError(err)
		qasmFile.RemoveFile()
		return err
	}

	logger.LogAction("Adding job to db")
	err = server.Database.AddJob(jobProperties, qasmFilePath, jobId)
	if err != nil {
		logger.LogError(err)
		qasmFile.RemoveFile()
		return err
	}

	logger.LogAction("Adding job to queue")
	err = server.QueueChannel.AddJob(server.QueueName, jobId)
	if err != nil {
		logger.LogError(err)
		qasmFile.RemoveFile()
		return err
	}

	return request.SendAndClose(&jobsServerProto.PendingJob{Id: jobId})
}
