package gmail

import (
	"context"
	"encoding/base64"
	"golang.org/x/oauth2"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/nonedotone/golog"
	"github.com/nonedotone/smtp-proxy/config"
	"github.com/nonedotone/smtp-proxy/structs"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Gmail struct {
	service *gmail.Service
}

func NewService(token *config.GmailToken, credentialsPath string) (*Gmail, error) {
	credentials, err := ReadGmailCredentialsOrDefault(credentialsPath)
	cfg, err := google.ConfigFromJSON(credentials, token.Permission)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, config.HttpClient)
	service, err := gmail.NewService(ctx,
		option.WithTokenSource(cfg.TokenSource(ctx, token.Oauth)))
	if err != nil {
		return nil, err
	}
	return &Gmail{service: service}, nil
}

func (m *Gmail) List(name, query string) ([]structs.Message, error) {
	return m.list(name, query)
}
func (m *Gmail) Send(msg []byte) error {
	return m.send(msg)
}

func (m *Gmail) list(name, query string) (msgs []structs.Message, err error) {
	golog.Debugf("list name %s, query %s\n", name, query)
	var (
		msgResp   *gmail.ListMessagesResponse
		pageToken string
		b         = true
	)
	for b || pageToken != "" {
		golog.Debugf("list param b %v, pageToken %s\n", b, pageToken)
		b = false
		//list
		if query != "" {
			msgResp, err = m.service.Users.Messages.List(name).Q(query).PageToken(pageToken).Do()
		} else {
			msgResp, err = m.service.Users.Messages.List(name).PageToken(pageToken).Do()
		}
		if err != nil {
			golog.Errorf("unable to retrieve mails: %v\n", err)
			return
		}
		if len(msgResp.Messages) == 0 {
			return
		}
		for _, msg := range msgResp.Messages {
			get, err := m.get(name, msg.Id)
			if err != nil {
				golog.Errorf("get name %s id %s error %v", name, msg.Id, err)
				continue
			}
			msgs = append(msgs, structs.Message{From: get.From, To: get.To, Subject: get.Subject})
		}
		pageToken = msgResp.NextPageToken
	}
	return
}
func (m *Gmail) get(name, msgId string) (parsemail.Email, error) {
	golog.Debugf("get user %s email %s\n", name, msgId)
	message, err := m.service.Users.Messages.Get(name, msgId).Format("raw").Do()
	if err != nil {
		golog.Errorf("get user %s message %s error %v\n", name, msgId, err)
		return parsemail.Email{}, err
	}
	rawBytes, err := base64.URLEncoding.DecodeString(message.Raw)
	if err != nil {
		golog.Errorf("msg id %s decode raw error %v\n", msgId, err)
		return parsemail.Email{}, err
	}
	email, err := parsemail.Parse(strings.NewReader(string(rawBytes)))
	if err != nil {
		golog.Errorf("msg id %s parse message %s error %v\n", msgId, string(rawBytes), err)
		return parsemail.Email{}, err
	}
	return email, nil
}
func (m *Gmail) send(rawBytes []byte) error {
	message := &gmail.Message{Raw: base64.URLEncoding.EncodeToString(rawBytes)}
	_, err := m.service.Users.Messages.Send("me", message).Do()
	return err
}
func (m *Gmail) sendMsg(msg string) error {
	message := &gmail.Message{Raw: base64.URLEncoding.EncodeToString([]byte(msg))}
	_, err := m.service.Users.Messages.Send("me", message).Do()
	return err
}
