package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/fdemchenko/arcus/internal/config"
	"gopkg.in/gomail.v2"
)

type MailSender struct {
	dialer         *gomail.Dialer
	sender         string
	templatesFS    fs.FS
	templatesCache map[string]*template.Template
}

func NewMailSender(cfg config.SMTPConfig, templatesFS fs.FS) *MailSender {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	cache := make(map[string]*template.Template)

	return &MailSender{dialer: dialer, sender: cfg.SenderAddress, templatesFS: templatesFS, templatesCache: cache}
}

func (m *MailSender) prepareTemplate(templateName string) (*template.Template, error) {
	if cachedTemplate, exists := m.templatesCache[templateName]; exists {
		return cachedTemplate, nil
	}
	parsedTemplate, err := template.ParseFS(m.templatesFS, templateName)
	if err != nil {
		return nil, fmt.Errorf("failer to parse email templates: %w", err)
	}
	return parsedTemplate, nil
}

func (m *MailSender) Send(to string, templateName string, data interface{}) error {
	t, err := m.prepareTemplate(templateName)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = t.ExecuteTemplate(subject, "subject", nil)
	if err != nil {
		return err
	}

	textBody := new(bytes.Buffer)
	err = t.ExecuteTemplate(textBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = t.ExecuteTemplate(htmlBody, "htmlBody", data)
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
