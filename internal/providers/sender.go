package providers

import (
	"context"
	"strings"
	"time"

	mailersend "github.com/mailersend/mailersend-go"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

type EmailConfig struct {
	ApiToken       string
	Subject        string
	Body           string
	FromEmail      string
	RecipientEmail string
}

func SendEmailMailerSender(config EmailConfig) (*string, error) {
	//mailersender provider can be validated with a dummy email to himself
	ms := mailersend.NewMailersend(config.ApiToken)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Email: config.FromEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Email: config.RecipientEmail,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(config.Subject)
	message.SetText(config.Body)

	res, error := ms.Email.Send(ctx, message)

	if error != nil {
		return nil, error
	}
	messageId := res.Header.Get("X-Message-Id")
	return &messageId, nil

}

func SendEmailMailgun(config EmailConfig) (*string, *string, error) {

	domain := *removeDomain(config.FromEmail)

	mg := mailgun.NewMailgun(domain, config.ApiToken)

	message := mg.NewMessage(
		config.FromEmail,
		config.Subject,
		config.Body,
		config.RecipientEmail,
	)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	resp, id, error := mg.Send(ctx, message)

	if error != nil {
		return nil, nil, error
	}

	return &resp, &id, nil
}

func validateDomainMailGun(apiKey, fromEmail string) (*string, error) {
	//paid feature
	v := mailgun.NewEmailValidator(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := v.ValidateEmail(ctx, fromEmail, false)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

func removeDomain(email string) *string {
	components := strings.Split(email, "@")
	domain := components[1]
	return &domain
}
