package onelogin

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/credentials"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/tokens"
)

// CacheDir is credentials cache dir
var CacheDir string

// Config provides configuration for API Clients
type Config struct {
	Endpoint     string
	ClientToken  string
	ClientSecret string
	Credentials  *credentials.Credentials
}

// NewConfig returns a new Config pointer
func NewConfig(endpoint string, clientToken string, clientSecret string) *Config {
	var c credentials.Value
	var v *credentials.Value
	if _, err := toml.DecodeFile(cacheFile(clientToken), &c); err == nil {
		v = &c
	}
	t := tokens.NewTokens()
	t.Endpoint = endpoint
	t.ClientToken = clientToken
	t.ClientSecret = clientSecret
	return &Config{
		Endpoint:     endpoint,
		ClientToken:  clientToken,
		ClientSecret: clientSecret,
		Credentials:  credentials.New(t, v),
	}
}

// Refresh load new credentials if necessary
func (c *Config) Refresh() error {
	return c.Credentials.Refresh()
}

// Save seves credentials value
func (c *Config) Save() error {
	if CacheDir != "" {
		fd, err := os.Create(cacheFile(c.ClientToken))
		if err != nil {
			return err
		}
		defer fd.Close()
		encoder := toml.NewEncoder(fd)
		creds, err := c.Credentials.Get()
		if err != nil {
			return err
		}
		return encoder.Encode(&creds)
	}
	return nil
}

func cacheFile(clientToken string) string {
	return path.Join(CacheDir, fmt.Sprintf("onelogin.%s.cache", clientToken))
}
