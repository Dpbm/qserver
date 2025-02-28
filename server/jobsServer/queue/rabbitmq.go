package queue

import (
	"context"
	"fmt"
	"time"

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

func (channel *RabbitMQChannel) AddJob(queueName string, jobId string) error {
	_, err := channel.Channel.QueueDeclare(
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

	err = channel.Channel.PublishWithContext(
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

	return err

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
