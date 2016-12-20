package alert

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/tanji/overseer/config"
)

type Alert struct {
	config.Alert
	Origin string
	Type   string
	Value  string
}

func (a *Alert) Email(auth config.SMTP) error {
	e := email.NewEmail()
	e.From = a.Sender
	e.To = a.Emails
	subj := fmt.Sprintf("Overseer | Alert | %s - %s: %s", a.Origin, a.Type, a.Value)
	e.Subject = subj
	text := fmt.Sprintf("Server: %s\nType:%s\nValue:%s", a.Origin, a.Type, a.Value)
	e.Text = []byte(text)
	var smtpauth smtp.Auth
	if auth.Username != "" && auth.Password != "" && auth.Server != "" {
		smtpauth = smtp.PlainAuth("", auth.Username, auth.Password, auth.Server)
	}
	err := e.Send(a.Destination, smtpauth)
	return err
}
