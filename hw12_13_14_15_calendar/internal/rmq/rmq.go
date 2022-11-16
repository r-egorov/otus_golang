package rmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQ struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	URI       string
	QueueName string
}

func New(uri, queueName string) *RMQ {
	return &RMQ{
		URI:       uri,
		QueueName: queueName,
	}
}

func (r *RMQ) Connect() error {
	conn, err := amqp.Dial(r.URI)
	if err != nil {
		return err
	}
	r.Conn = conn

	amqpCh, err := conn.Channel()
	if err != nil {
		return err
	}
	r.Channel = amqpCh

	_, err = amqpCh.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *RMQ) Close() error {
	err := r.Conn.Close()
	if err != nil {
		return err
	}
	err = r.Channel.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *RMQ) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := r.Channel.Consume(
		r.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r *RMQ) Publish(ctx context.Context, msg []byte) error {
	err := r.Channel.PublishWithContext(
		ctx,
		"",
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
