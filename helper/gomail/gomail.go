package gomail

import (
	"BE-hi-SPEC/config"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func Gomail(email string) error {
	config := config.InitConfig()
	sender := NewGmailSender(config.GOMAIL_NAME, config.GOMAIL_EMAIL, config.GOMAIL_PASS)

	subject := "Transaction PDF"
	content := `
	<h1>hi'SPEC</h1>
	<pre>
	Here is the copy of your transaction pdf file.
	
	Best regards, hi'SPEC Admin.
	</pre>
	`

	to := []string{email}
	attachFiles := []string{"helper/gofpdf/invoice.pdf"}

	err := sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	if err != nil {
		return err
	}

	return nil
}
