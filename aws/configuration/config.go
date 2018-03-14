package configuration

import (
	"fmt"
	"os"
	"path"

	"github.com/go-ini/ini"
)

// Config represents ~/.aws/config handler
type Config struct {
	file    string
	profile string
}

// NewConfig creates a Config
func NewConfig(dir string, profile string) *Config {
	return &Config{
		file:    path.Join(dir, "config"),
		profile: profile,
	}
}

// Save to ~/.aws/config
func (c *Config) Save(region string) error {
	configIni, err := ini.Load(c.file)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		configIni = ini.Empty()
	}
	section := configIni.Section(fmt.Sprintf("profile %s", c.profile))
	k, err := section.GetKey("region")
	if err != nil {
		_, err := section.NewKey("region", region)
		if err != nil {
			return err
		}
	} else {
		k.SetValue(region)
	}
	return configIni.SaveTo(c.file)
}
