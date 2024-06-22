package providers

import (
	"context"
	"strings"
	"time"

	mailersend "github.com/mailersend/mailersend-go"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

func SendEmailMailerSender(apiToken, subject, text, fromEmail, recipientEmail string) (*string, error) {

	//ms := mailersend.NewMailersend(os.Getenv("MAILERSEND_API_KEY"))
	ms := mailersend.NewMailersend(apiToken)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Email: fromEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Email: recipientEmail,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetText(text)

	res, error := ms.Email.Send(ctx, message)

	if error != nil {
		return nil, error
	}
	messageId := res.Header.Get("X-Message-Id")
	return &messageId, nil

}

func sendEmailMailgun(apiToken, subject, text, fromEmail, recipientEmail string) (*string, *string, error) {

	domain := *removeDomain(fromEmail)

	mg := mailgun.NewMailgun(domain, apiToken)

	message := mg.NewMessage(
		fromEmail,
		subject,
		text,
		recipientEmail,
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
