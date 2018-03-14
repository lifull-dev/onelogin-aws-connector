package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Config stores config
type Config struct {
	Service map[string]*ServiceConfig `toml:"service"`
	App     map[string]*AppConfig     `toml:"app"`
	file    string                    `toml:"-"`
}

// ServiceConfig stores initialized data
type ServiceConfig struct {
	Endpoint        string `toml:"endpoint"`
	ClientToken     string `toml:"client_token"`
	ClientSecret    string `toml:"client_secret"`
	Subdomain       string `toml:"subdomain"`
	UsernameOrEmail string `toml:"username_or_email"`
}

// AppConfig stores configured data
type AppConfig struct {
	AppID        string `toml:"app_id"`
	RoleArn      string `toml:"role_arn"`
	PrincipalArn string `toml:"principal_arn"`
}

// Load creates a Loaded Config
func Load(file string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(file, &config); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		config = Config{
			Service: map[string]*ServiceConfig{},
			App:     map[string]*AppConfig{},
		}
	}
	if config.App == nil {
		config.App = map[string]*AppConfig{}
	}
	config.file = file
	return &config, nil
}

// Save to persistent store
func (c Config) Save() error {
	fd, err := os.Create(c.file)
	if err != nil {
		return err
	}
	defer fd.Close()
	encoder := toml.NewEncoder(fd)
	return encoder.Encode(c)
}
