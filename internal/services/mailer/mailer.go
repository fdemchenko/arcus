package mailer

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/fdemchenko/arcus/internal/config"
	"github.com/fdemchenko/arcus/templates"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer          *gomail.Dialer
	welcomeTemplate *template.Template
	sender          string
}

func New(cfg config.SMTPConfig) (*Mailer, error) {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	t, err := template.ParseFS(templates.TemplatesFS, "user_welcome.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failer to parse email templates: %w", err)
	}

	return &Mailer{dialer: dialer, sender: cfg.SenderAddress, welcomeTemplate: t}, nil
}

func (m *Mailer) Send(to string, data interface{}) error {
	subject := new(bytes.Buffer)
	err := m.welcomeTemplate.ExecuteTemplate(subject, "subject", nil)
	if err != nil {
		return err
	}

	textBody := new(bytes.Buffer)
	err = m.welcomeTemplate.ExecuteTemplate(textBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = m.welcomeTemplate.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.sender)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/plain", textBody.String())
	message.AddAlternative("text/html", htmlBody.String())

	return m.dialer.DialAndSend(message)
}
