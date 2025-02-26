package types

import (
	"github.com/Dpbm/jobsServer/db"
	jobsServerProto "github.com/Dpbm/jobsServer/proto"
	"github.com/Dpbm/jobsServer/queue"
)

type JobsServer struct {
	jobsServerProto.UnimplementedJobsServer
	QueueChannel *queue.QueueChannel
	Database     *db.DB
}
