package mail

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"html/template"

	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

//go:embed "templates"
var templateFS embed.FS

type SendParams struct {
	From    string
	To      string
	Subject string
	File    string
	Data    any
}

type Config struct {
	APIKey     string
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
}

type Mail struct {
	ctx    context.Context
	client *resend.Client
	config Config
}

func New(ctx context.Context, config Config) *Mail {
	return &Mail{
		ctx:    ctx,
		client: resend.NewClient(config.APIKey),
		config: config,
	}
}

func (m *Mail) Send(p SendParams) error {
	tmplLocation := fmt.Sprintf("templates/%s", p.File)

	tmpl, err := template.New("email").ParseFS(templateFS, tmplLocation)
	if err != nil {
		logging.Error("error on parse template", err,
			zap.String("journey", "mail"))
		return fault.New(
			"failed to parse template",
			fault.WithHTTPCode(http.StatusInternalServerError),
			fault.WithTag(fault.MAILER_ERROR),
			fault.WithError(err),
		)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, p.Data); err != nil {
		logging.Error("error on execute template", err,
			zap.String("journey", "mail"))
		return fault.New(
			"failed to execute template",
			fault.WithHTTPCode(http.StatusInternalServerError),
			fault.WithTag(fault.MAILER_ERROR),
			fault.WithError(err),
		)
	}

	params := &resend.SendEmailRequest{
		From:    p.From,
		To:      []string{p.To},
		Html:    body.String(),
		Subject: p.Subject,
	}

	return m.send(params, m.config.MaxRetries)
}

func (m *Mail) send(p *resend.SendEmailRequest, maxRetries int) error {
	var mailerErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(m.ctx, m.config.Timeout)
		defer cancel()

		_, err := m.client.Emails.SendWithContext(ctx, p)
		if err == nil {
			logging.Info("email sent",
				zap.Int("attempt", attempt),
				zap.String("to", p.To[0]),
			)
			return nil
		}

		logging.Error("error on send email", err,
			zap.Int("attempt", attempt),
			zap.String("to", p.To[0]),
		)

		mailerErr = err
		time.Sleep(m.config.RetryDelay)
	}

	return fault.New(
		fmt.Sprintf("error on send email after %d attemps", maxRetries),
		fault.WithHTTPCode(http.StatusInternalServerError),
		fault.WithTag(fault.MAILER_ERROR),
		fault.WithError(mailerErr),
	)
}
