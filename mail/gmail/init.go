package gmail

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nonedotone/golog"
	spcfg "github.com/nonedotone/smtp-proxy/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const (
	MailPermissionReadonly = gmail.GmailReadonlyScope
	MailGmailSendScope     = gmail.GmailSendScope
)

const (
	MailCredentials = `{"installed":{"client_id":"744825144845-l91ulno2ga9t9g7k6oqaedpqjr22r3s7.apps.googleusercontent.com","project_id":"smtp-proxy-356602","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"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","redirect_uris":["http://localhost"]}}`
)

func ReadGmailCredentialsOrDefault(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		golog.Warn("use default gmail credentials")
		return []byte(MailCredentials), nil
	}
	golog.Infof("read gmail credentials from path %s\n", path)
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s error %v", path, err)
	}
	return bz, nil
}

func GetTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	fmt.Print("auth code -> ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, spcfg.HttpClient)
	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		golog.Debugf("config exchange error %v", err)
		return nil, err
	}
	return tok, nil
}
func InitConfig(permission string, credentials []byte) (*spcfg.Config, error) {
	golog.Debugf("init config permission %s, credentials %s\n", permission, credentials)
	auth2Cfg, err := google.ConfigFromJSON(credentials, permission)
	if err != nil {
		golog.Debugf("google config from json error %v\n", err)
		return nil, err
	}
	oauth, err := GetTokenFromWeb(auth2Cfg)
	if err != nil {
		golog.Debugf("google token from web error %v\n", err)
		return nil, err
	}
	gmailToken := &spcfg.GmailToken{
		Permission: permission,
		Oauth:      oauth,
	}
	return &spcfg.Config{Type: spcfg.GmailType, Token: gmailToken}, nil
}
