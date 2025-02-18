package mail

import (
	"context"
	"errors"

	"github.com/resend/resend-go/v2"
)

var (
	client *resend.Client
)

func InitMailer(apiKey string) {
	client = resend.NewClient(apiKey)
}

type EmailParams struct {
	From    string
	To      string
	Subject string
	Html    string
}

type EmailResponse struct {
	ID  string
	Err error
}

func Send(params EmailParams) (string, error) {
	ctx := context.TODO()

	to := params.To
	if to == "" {
		return "", errors.New("to is required")
	}

	html := params.Html
	if html == "" {
		return "", errors.New("html is required")
	}

	from := params.From
	if from == "" {
		from = "Meu Buzufba <no-reply@condosnap.com.br>"
	}

	subject := params.Subject
	if subject == "" {
		subject = "Meu Buzufba: Tem email novo pra"
	}

	email := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    params.Html,
	}

	result, err := client.Emails.SendWithContext(ctx, email)
	if err != nil {
		return "", err
	}

	return result.Id, nil
}
