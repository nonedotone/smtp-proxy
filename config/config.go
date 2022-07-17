package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	GmailType = "gmail"
)

type Config struct {
	Type  string      `json:"type"`
	Token interface{} `json:"token"`
}

type GmailToken struct {
	Permission string        `json:"permission"`
	Oauth      *oauth2.Token `json:"oauth"`
}

func (cfg *Config) UnmarshalJSON(data []byte) error {
	type cloneType Config
	rawMsg := json.RawMessage{}
	cfg.Token = &rawMsg

	if err := json.Unmarshal(data, (*cloneType)(cfg)); err != nil {
		return err
	}

	switch cfg.Type {
	case GmailType:
		params := new(GmailToken)
		if err := json.Unmarshal(rawMsg, params); err != nil {
			return err
		}
		cfg.Token = params
	default:
		return errors.New("nonsupport type")
	}
	return nil
}

func LoadConfig(f string) (*Config, error) {
	bz, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	err = json.Unmarshal(bz, &cfg)
	return &cfg, err
}

var HttpClient = &http.Client{Transport: http.DefaultTransport, Timeout: 60 * time.Second}
