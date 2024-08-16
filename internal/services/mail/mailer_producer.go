package mail

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MailerProducer struct {
	channel *amqp.Channel
}

func NewMailerProducer(channel *amqp.Channel) (*MailerProducer, error) {
	_, err := channel.QueueDeclare(
		EmailQueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &MailerProducer{
		channel: channel,
	}, nil
}

func (mp *MailerProducer) Publish(command SendEmailCommand[any]) error {
	js, err := json.Marshal(command)
	if err != nil {
		return err
	}
	return mp.channel.Publish(
		"",
		EmailQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         js,
		},
	)
}
