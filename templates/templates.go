package templates

import "embed"

//go:embed *.tmpl
var TemplatesFS embed.FS

type UserWelcomeData struct {
	Token string
	Host  string
}
