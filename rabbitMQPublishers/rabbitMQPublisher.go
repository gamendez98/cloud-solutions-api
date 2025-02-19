package rabbitMQPublishers

import (
	"cloud-solutions-api/config"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	needsDeclare bool
	exchangeName string
	queueName    string
}

func (publisher *RabbitMQPublisher) PublishJSON(message interface{}) error {
	if publisher.ch == nil {
		return fmt.Errorf("channel is not initialized")
	}

	marshalled, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := publisher.ch.Publish(
		publisher.exchangeName, // exchange
		"",                     // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalled,
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func NewRabbitMQPublisher(exchangeName string, queueName string) (*RabbitMQPublisher, error) {
	publisher := &RabbitMQPublisher{
		exchangeName: exchangeName,
		queueName:    queueName,
	}
	if err := publisher.Connect(); err != nil {
		return nil, err
	}
	if err := publisher.DeclareAndBindQueue(); err != nil {
		return nil, err
	}
	return publisher, nil
}

func (publisher *RabbitMQPublisher) Connect() error {
	conn, err := GetConnection()
	if err != nil {
		return err
	}
	publisher.conn = conn
	publisher.ch, err = conn.Channel()
	if err != nil {
		return err
	}
	return nil
}

func (publisher *RabbitMQPublisher) DeclareAndBindQueue() error {

	if err := publisher.ch.ExchangeDeclare(
		publisher.exchangeName, // exchange name
		"direct",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := publisher.ch.QueueDeclare(
		publisher.queueName, // queue name
		true,                // durable
		false,               // auto-deleted
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := publisher.ch.QueueBind(
		queue.Name,             // queue name
		"",                     // routing key
		publisher.exchangeName, // exchange name
		false,                  // no-wait
		nil,                    // arguments
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

var connection *amqp.Connection

func GetConnection() (*amqp.Connection, error) {
	if connection != nil && !connection.IsClosed() {
		return connection, nil
	}
	conf := config.GetConfig()
	uri :=
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			conf.RabbitMQUsername,
			conf.RabbitMQPassword,
			conf.RabbitMQHost,
			conf.RabbitMQPort)
	return amqp.Dial(uri)
}
