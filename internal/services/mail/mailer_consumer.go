package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

const EmailQueueName = "welcome_emails"

type Sender interface {
	Send(to string, templateName string, data interface{}) error
}

type MailerConsumer struct {
	sender  Sender
	channel *amqp.Channel
	logger  *slog.Logger
}

func NewMailerConsumer(sender Sender, channel *amqp.Channel, logger *slog.Logger) (*MailerConsumer, error) {
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
		logger:  logger,
	}, nil
}

func (mc *MailerConsumer) StartConsuming() error {
	const op = "mail.MailerConsumer.StartConsuming"
	logger := mc.logger.With(slog.String("op", op))
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
			err = mc.handleDelivery(delivery)
			if err != nil {
				logger.Error("failed to handle amqp delivery", slog.String("error", err.Error()))
			}
			err := delivery.Ack(false)
			if err != nil {
				logger.Error("failed to ack delivery", slog.String("error", err.Error()))
			}
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

	err = mc.sender.Send(sendEmailCommand.To, sendEmailCommand.TemplateName, sendEmailCommand.TemplateData)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
