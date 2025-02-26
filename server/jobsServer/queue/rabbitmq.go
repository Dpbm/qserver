package queue

import (
	"fmt"

	logger "github.com/Dpbm/shared/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct{}

type RabbitMQConnection struct {
	Connection *amqp.Connection
}

type RabbitMQChannel struct {
	Channel *amqp.Channel
}

func (channel *RabbitMQChannel) Close() {
	channel.Channel.Close()
}

func (connection *RabbitMQConnection) Close() {
	connection.Connection.Close()
}

func (connection *RabbitMQConnection) CreateChannel() QueueChannel {
	channel, err := connection.Connection.Channel()

	if err != nil {
		logger.LogFatal(err) // it will exit with status 1
	}

	return &RabbitMQChannel{Channel: channel}
}

func (queue *RabbitMQ) ConnectQueue(host string, port int, username string, password string) QueueConnection {
	rabbitmqServerUrl := fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, host, port)
	connection, err := amqp.Dial(rabbitmqServerUrl)

	if err != nil {
		logger.LogFatal(err) // it will exit with status 1
	}

	return &RabbitMQConnection{Connection: connection}
}
