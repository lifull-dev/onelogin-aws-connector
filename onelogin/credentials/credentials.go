package credentials

import (
	"time"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin/tokens"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/tokens/tokensiface"
)

// Credentials provides credentials for API Clients
type Credentials struct {
	Credentials *Value
	Tokens      tokensiface.TokensAPI
}

// Value provides credentials for API Clients
type Value struct {
	AccessToken      string
	RefreshToken     string
	CreatedAt        time.Time
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

// New returns a new Credentials pointer
func New(t tokensiface.TokensAPI, c *Value) *Credentials {
	return &Credentials{
		Credentials: c,
		Tokens:      t,
	}
}

// Get returns the credentials value, or error
func (c *Credentials) Get() (Value, error) {
	if err := c.Refresh(); err != nil {
		return Value{}, err
	}
	return *c.Credentials, nil
}

// Refresh load new credentials if necessary
func (c *Credentials) Refresh() error {
	var res *tokens.GenerateResponse
	var err error
	if c.Credentials != nil {
		creds := c.Credentials
		if creds.availavle() {
			return nil
		}
		if creds.refreshable() {
			input := &tokens.RefreshRequest{
				AccessToken:  c.Credentials.AccessToken,
				RefreshToken: c.Credentials.RefreshToken,
			}
			res, err = c.Tokens.Refresh(input)
			if err != nil {
				if err.Error() != "[401] Unauthorized: Invalid Token" {
					return err
				}
				res, err = c.Tokens.Generate()
				if err != nil {
					return err
				}
			}
		} else {
			res, err = c.Tokens.Generate()
			if err != nil {
				return err
			}
		}
	} else {
		res, err = c.Tokens.Generate()
		if err != nil {
			return err
		}
	}
	createdAt, err := time.Parse("2006-01-02T15:04:05Z", res.CreatedAt)
	if err != nil {
		return err
	}
	createdAt = createdAt.UTC()
	accessExpiresAt := createdAt.Add(time.Duration(res.ExpiresIn * 1000000000))
	refreshExpiresAt := createdAt.Add(45 * 24 * time.Hour)

	c.Credentials = &Value{
		AccessToken:      res.AccessToken,
		RefreshToken:     res.RefreshToken,
		CreatedAt:        createdAt,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}
	return nil
}

func (c *Value) availavle() bool {
	return time.Now().Before(c.AccessExpiresAt)
}

func (c *Value) refreshable() bool {
	return time.Now().Before(c.RefreshExpiresAt)
}
