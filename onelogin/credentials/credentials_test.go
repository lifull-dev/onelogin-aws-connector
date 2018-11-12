package credentials

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin/tokens"
)

func TestNew(t *testing.T) {
	creds := &Value{}
	tokenapi := tokens.NewTokens()
	c := New(tokenapi, creds)
	if c.Credentials != creds {
		t.Errorf("Value = %v, want %v", c.Credentials, creds)
	}
	if c.Tokens != tokenapi {
		t.Errorf("Tokens = %v, want %v", c.Tokens, tokenapi)
	}
}

type TokenAPIMock struct {
	GenerateResponse       *tokens.GenerateResponse
	RefreshResponse        *tokens.RefreshResponse
	RefreshRequestVerifier func(*tokens.RefreshRequest) error
	Error                  error
}

func (t *TokenAPIMock) Generate() (*tokens.GenerateResponse, error) {
	return t.GenerateResponse, t.Error
}

func (t *TokenAPIMock) Refresh(input *tokens.RefreshRequest) (*tokens.RefreshResponse, error) {
	if err := t.RefreshRequestVerifier(input); err != nil {
		return nil, err
	}
	return t.RefreshResponse, t.Error
}

func TestCredentialsGet(t *testing.T) {
	t.Run("when Refresh() success", func(t *testing.T) {
		n := time.Now().UTC()
		v := &Value{
			CreatedAt:        n,
			AccessExpiresAt:  n.Add(10 * time.Second),
			RefreshExpiresAt: n.Add(100 * time.Second),
		}
		a := &TokenAPIMock{}
		c := &Credentials{
			Credentials: v,
			Tokens:      a,
		}
		got, err := c.Get()
		if err != nil {
			t.Errorf("Credentials.Get() error = %v", err)
		}
		if !reflect.DeepEqual(got, *v) {
			t.Errorf("Credentials.Get() = %v, want %v", got, v)
		}
	})
	t.Run("when Refresh() error", func(t *testing.T) {
		e := fmt.Errorf("error")
		a := &TokenAPIMock{
			Error: e,
		}
		c := &Credentials{
			Credentials: nil,
			Tokens:      a,
		}
		_, err := c.Get()
		if err != e {
			t.Errorf("Credentials.Get() error = %v", err)
		}
	})
}

func TestCredentialsRefresh(t *testing.T) {
	t.Run("when No Credentials", func(t *testing.T) {
		n, _ := time.Parse("2006-01-02T15:04:05Z", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
		n = n.UTC()
		a := &TokenAPIMock{
			GenerateResponse: &tokens.GenerateResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				CreatedAt:    n.Format("2006-01-02T15:04:05Z"),
				ExpiresIn:    100,
			},
		}
		expected := Value{
			AccessToken:      "access-token",
			RefreshToken:     "refresh-token",
			CreatedAt:        n,
			AccessExpiresAt:  n.Add(100 * time.Second),
			RefreshExpiresAt: n.Add(45 * 24 * time.Hour),
		}
		c := &Credentials{
			Credentials: nil,
			Tokens:      a,
		}
		got, err := c.Get()
		if err != nil {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Credentials.Get() = %#v, want %#v", got, expected)
		}
	})
	t.Run("when No Credentials Error", func(t *testing.T) {
		e := fmt.Errorf("error")
		a := &TokenAPIMock{
			Error: e,
		}
		c := &Credentials{
			Credentials: nil,
			Tokens:      a,
		}
		_, err := c.Get()
		if err != e {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
	})
	t.Run("when available Credentials", func(t *testing.T) {
		n, _ := time.Parse("2006-01-02T15:04:05Z", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
		n = n.UTC()
		a := &TokenAPIMock{}
		v := &Value{
			AccessToken:      "access-token",
			RefreshToken:     "refresh-token",
			CreatedAt:        n.Add(-100 * time.Second),
			AccessExpiresAt:  n.Add(100 * time.Second),
			RefreshExpiresAt: n.Add(45 * 24 * time.Hour),
		}
		expected := *v
		c := &Credentials{
			Credentials: v,
			Tokens:      a,
		}
		got, err := c.Get()
		if err != nil {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Credentials.Get() = %#v, want %#v", got, expected)
		}
	})
	t.Run("when unavailable and refreshable Credentials", func(t *testing.T) {
		n, _ := time.Parse("2006-01-02T15:04:05Z", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
		n = n.UTC()
		a := &TokenAPIMock{
			RefreshResponse: &tokens.RefreshResponse{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				CreatedAt:    n.Format("2006-01-02T15:04:05Z"),
				ExpiresIn:    100,
			},
			RefreshRequestVerifier: func(input *tokens.RefreshRequest) error {
				return nil
			},
		}
		v := &Value{
			AccessToken:      "access-token",
			RefreshToken:     "refresh-token",
			CreatedAt:        n.Add(-11 * time.Hour),
			AccessExpiresAt:  n.Add(-1 * time.Hour),
			RefreshExpiresAt: n.Add(44 * 24 * time.Hour),
		}
		expected := Value{
			AccessToken:      "new-access-token",
			RefreshToken:     "new-refresh-token",
			CreatedAt:        n,
			AccessExpiresAt:  n.Add(100 * time.Second),
			RefreshExpiresAt: n.Add(45 * 24 * time.Hour),
		}
		c := &Credentials{
			Credentials: v,
			Tokens:      a,
		}
		got, err := c.Get()
		if err != nil {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Credentials.Get() = %v, want %v", got, expected)
		}
	})
	t.Run("when unavailable and refreshable Credentials Error", func(t *testing.T) {
		n, _ := time.Parse("2006-01-02T15:04:05Z", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
		n = n.UTC()
		e := fmt.Errorf("error")
		a := &TokenAPIMock{
			RefreshRequestVerifier: func(input *tokens.RefreshRequest) error {
				return nil
			},
			Error: e,
		}
		v := &Value{
			AccessToken:      "access-token",
			RefreshToken:     "refresh-token",
			CreatedAt:        n.Add(-11 * time.Hour),
			AccessExpiresAt:  n.Add(-1 * time.Hour),
			RefreshExpiresAt: n.Add(44 * 24 * time.Hour),
		}
		c := &Credentials{
			Credentials: v,
			Tokens:      a,
		}
		_, err := c.Get()
		if err != e {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
	})
	t.Run("when unrefreshable Credentials", func(t *testing.T) {
		n, _ := time.Parse("2006-01-02T15:04:05Z", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
		n = n.UTC()
		a := &TokenAPIMock{
			GenerateResponse: &tokens.GenerateResponse{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				CreatedAt:    n.Format("2006-01-02T15:04:05Z"),
				ExpiresIn:    100,
			},
			RefreshRequestVerifier: func(input *tokens.RefreshRequest) error {
				return nil
			},
		}
		v := &Value{
			AccessToken:      "access-token",
			RefreshToken:     "refresh-token",
			CreatedAt:        n.Add(-45 * 24 * time.Hour),
			AccessExpiresAt:  n.Add(-44 * 24 * time.Hour),
			RefreshExpiresAt: n.Add(-1 * 24 * time.Hour),
		}
		expected := Value{
			AccessToken:      "new-access-token",
			RefreshToken:     "new-refresh-token",
			CreatedAt:        n,
			AccessExpiresAt:  n.Add(100 * time.Second),
			RefreshExpiresAt: n.Add(45 * 24 * time.Hour),
		}
		c := &Credentials{
			Credentials: v,
			Tokens:      a,
		}
		got, err := c.Get()
		if err != nil {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Credentials.Get() = %v, want %v", got, expected)
		}
	})
	t.Run("when unrefreshable Credentials Error", func(t *testing.T) {
		n, _ := time.Parse("2006-01-02T15:04:05Z", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
		n = n.UTC()
		e := fmt.Errorf("error")
		a := &TokenAPIMock{
			Error: e,
		}
		v := &Value{
			AccessToken:      "access-token",
			RefreshToken:     "refresh-token",
			CreatedAt:        n.Add(-45 * 24 * time.Hour),
			AccessExpiresAt:  n.Add(-44 * 24 * time.Hour),
			RefreshExpiresAt: n.Add(-1 * 24 * time.Hour),
		}
		c := &Credentials{
			Credentials: v,
			Tokens:      a,
		}
		_, err := c.Get()
		if err != e {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
	})
	t.Run("when invalid refresh token", func(t *testing.T) {
		e := fmt.Errorf("[401] Unauthorized: Invalid Token")
		n, _ := time.Parse("2006-01-02T15:04:05Z", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
		a := &TokenAPIMock{
			GenerateResponse: &tokens.GenerateResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				CreatedAt:    n.Format("2006-01-02T15:04:05Z"),
				ExpiresIn:    100,
			},
			RefreshRequestVerifier: func(t *tokens.RefreshRequest) error {
				return e
			},
		}
		v := &Value{
			AccessToken:      "access-token",
			RefreshToken:     "refresh-token",
			CreatedAt:        n,
			AccessExpiresAt:  n.Add(-10 * time.Second),
			RefreshExpiresAt: n.Add(100 * time.Second),
		}
		c := &Credentials{
			Credentials: v,
			Tokens:      a,
		}
		_, err := c.Get()
		if err != nil {
			t.Errorf("Credentials.Get() error = %#v", err)
		}
	})
}
