package mail

import (
	"bytes"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const EmailQueueName = "welcome_emails"

type Sender interface {
	Send(to string, templateName string, data interface{}) error
}

type MailerConsumer struct {
	sender  Sender
	channel *amqp.Channel
}

func NewMailerConsumer(sender Sender, channel *amqp.Channel) (*MailerConsumer, error) {
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

	return &MailerConsumer{
		sender:  sender,
		channel: channel,
	}, nil
}

func (mc *MailerConsumer) StartConsuming() error {
	deliveriesChan, err := mc.channel.Consume(
		EmailQueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveriesChan {
			_ = mc.handleDelivery(delivery)
			delivery.Ack(false)
		}
	}()
	return nil
}

func (mc *MailerConsumer) handleDelivery(delivery amqp.Delivery) error {
	var sendEmailCommand SendEmailCommand[UserWelcomeData]
	err := json.NewDecoder(bytes.NewReader(delivery.Body)).Decode(&sendEmailCommand)
	if err != nil {
		return fmt.Errorf("cannot decode json SendEmailCommand: %w", err)
	}

	return mc.sender.Send(sendEmailCommand.To, sendEmailCommand.TemplateName, sendEmailCommand.TemplateData)
}
