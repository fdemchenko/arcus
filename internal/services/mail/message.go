package mail

type SendEmailCommand[T any] struct {
	To           string `json:"to"`
	TemplateName string `json:"template_name"`
	TemplateData T      `json:"template_data"`
}

type UserWelcomeData struct {
	Token string `json:"token"`
}
