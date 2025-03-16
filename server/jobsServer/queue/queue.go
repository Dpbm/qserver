package queue

type QueueChannel interface {
	Close()
	AddJob(queueName string, jobId string) error
}

type QueueConnection interface {
	Close()
	CreateChannel() QueueChannel
}

type Queue interface {
	ConnectQueue(username string, password string, host string, port uint16) QueueConnection
}
