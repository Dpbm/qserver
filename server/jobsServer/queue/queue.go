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
	ConnectQueue(host string, port int, username string, password string) QueueConnection
}
