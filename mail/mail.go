package mail

import (
	"fmt"

	"github.com/nonedotone/smtp-proxy/config"
	"github.com/nonedotone/smtp-proxy/mail/gmail"
)

type Boxer interface {
	Send([]byte) error
}

type Mail struct {
	cfg *config.Config
	box Boxer
}

func NewMail(cfg *config.Config, gmailCredentials string) (*Mail, error) {
	var box Boxer
	var err error
	switch cfg.Type {
	case config.GmailType:
		t, ok := cfg.Token.(*config.GmailToken)
		if !ok {
			return nil, fmt.Errorf("token type error")
		}

		box, err = gmail.NewService(t, gmailCredentials)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("nonsupport type")
	}
	return &Mail{cfg: cfg, box: box}, nil
}

func (m *Mail) Send(from, to, subject, msg string) error {
	rawBytes := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, msg))
	return m.box.Send(rawBytes)
}
func (m *Mail) SendMessage(bz []byte) error {
	return m.box.Send(bz)
}
