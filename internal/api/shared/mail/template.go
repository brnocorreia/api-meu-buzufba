package mail

import (
	"bytes"
	"html/template"
	"path/filepath"
)

const (
	templatesPath = "internal/api/shared/mail/templates"
)

type WelcomeEmailData struct {
	Name            string
	VerificationURL string
}

func ParseWelcomeTemplate(data WelcomeEmailData) (string, error) {
	templatePath := filepath.Join(templatesPath, "welcome_user.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
