package configuration

import (
	"os"
	"path"

	"github.com/go-ini/ini"
)

// Credentials represents ~/.aws/credentials handler
type Credentials struct {
	file    string
	profile string
}

// NewCredentials creates a Credentials
func NewCredentials(dir string, profile string) *Credentials {
	return &Credentials{
		file:    path.Join(dir, "credentials"),
		profile: profile,
	}
}

// Save to ~/.aws/credentials
func (c *Credentials) Save(options map[string]string) error {
	credsIni, err := ini.Load(c.file)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		credsIni = ini.Empty()
	}
	section := credsIni.Section(c.profile)
	for key, value := range options {
		k, err := section.GetKey(key)
		if err != nil {
			_, err := section.NewKey(key, value)
			if err != nil {
				return err
			}
		} else {
			k.SetValue(value)
		}
	}
	return credsIni.SaveTo(c.file)
}
