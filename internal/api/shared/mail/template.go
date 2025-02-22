package mail

import (
	"bytes"
	"html/template"
	"path/filepath"
)

const (
	templatesPath = "internal/api/shared/mail/templates"
)

func ParseTemplate(templatePath string, data interface{}) (string, error) {
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

type WelcomeEmailData struct {
	Name            string
	VerificationURL string
}

func ParseWelcomeTemplate(data WelcomeEmailData) (string, error) {
	templatePath := filepath.Join(templatesPath, "welcome_user.html")
	return ParseTemplate(templatePath, data)
}

type VerifyEmailData struct {
	Name            string
	VerificationURL string
}

func ParseVerifyEmailTemplate(data VerifyEmailData) (string, error) {
	templatePath := filepath.Join(templatesPath, "verify_email.html")
	return ParseTemplate(templatePath, data)
}
