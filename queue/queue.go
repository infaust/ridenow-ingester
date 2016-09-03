package queue

import "github.com/streadway/amqp"

type QueueProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewQueueProducer(dataSourceName string) (*QueueProducer, error) {
	conn, err := amqp.Dial(dataSourceName)
	if err != nil {
		return nil, err
	}
	// defer conn.Close() // ?
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// defer ch.Close() // ?
	err = ch.ExchangeDeclare(
		"ridenow_matcher", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}
	return &QueueProducer{conn, ch}, nil
}

func (q *QueueProducer) Send(routing string, bytes []byte) error {
	err := q.channel.Publish(
		"ridenow_matcher", // exchange
		routing,           // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})
	if err != nil {
		return err
	}
	return nil
}
