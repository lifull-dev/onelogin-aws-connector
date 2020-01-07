package onelogin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin/credentials"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/tokens"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/tokens/tokensiface"
)

func TestNewConfigFileNotExists(t *testing.T) {
	CacheDir = os.TempDir()
	config := NewConfig("endpoint", "client-token", "client-secret")
	if config.Endpoint != "endpoint" {
		t.Errorf("%s is not equal %s", config.Endpoint, "endpoint")
	}
	if config.ClientToken != "client-token" {
		t.Errorf("%s is not equal %s", config.Endpoint, "client-token")
	}
	if config.ClientSecret != "client-secret" {
		t.Errorf("%s is not equal %s", config.Endpoint, "client-secret")
	}
	if config.Credentials.Credentials != nil {
		t.Errorf("%s is not nil", config.Credentials.Credentials)
	}
	if config.Credentials.Tokens == nil {
		t.Error("config.Credentials.Tokens is nil")
	}
}

func TestNewConfigFileExists(t *testing.T) {
	CacheDir = os.TempDir()

	now := time.Now()
	cache := fmt.Sprintf(`AccessToken = "access-token"
RefreshToken = "refresh-token"
CreatedAt = %s
AccessExpiresAt = %s
RefreshExpiresAt = %s`,
		now.Format("2006-01-02T15:04:05Z"),
		now.Add(2*time.Second).Format("2006-01-02T15:04:05Z"),
		now.Add(3*time.Second).Format("2006-01-02T15:04:05Z"))
	file := path.Join(CacheDir, fmt.Sprintf("onelogin.%s.cache", "client-token"))
	if err := ioutil.WriteFile(file, []byte(cache), 0666); err != nil {
		t.Errorf("%#v", err)
	}
	defer os.Remove(file)
	config := NewConfig("endpoint", "client-token", "client-secret")
	if config.Endpoint != "endpoint" {
		t.Errorf("%s is not equal %s", config.Endpoint, "endpoint")
	}
	if config.ClientToken != "client-token" {
		t.Errorf("%s is not equal %s", config.Endpoint, "client-token")
	}
	if config.ClientSecret != "client-secret" {
		t.Errorf("%s is not equal %s", config.Endpoint, "client-secret")
	}
	creds := config.Credentials.Credentials
	if creds.AccessToken != "access-token" {
		t.Errorf("%v is not equal %v", creds.AccessToken, "access-token")
	}
	if creds.RefreshToken != "refresh-token" {
		t.Errorf("%v is not equal %v", creds.RefreshToken, "refresh-token")
	}
	if creds.CreatedAt.Format("2006-01-02T15:04:05Z") != now.Format("2006-01-02T15:04:05Z") {
		t.Errorf("%v is not equal %v", creds.CreatedAt, now.Format("2006-01-02T15:04:05Z"))
	}
	if creds.AccessExpiresAt.Format("2006-01-02T15:04:05Z") != now.Add(2*time.Second).Format("2006-01-02T15:04:05Z") {
		t.Errorf("%v is not equal %v", creds.AccessExpiresAt, now.Add(2*time.Second).Format("2006-01-02T15:04:05Z"))
	}
	if creds.RefreshExpiresAt.Format("2006-01-02T15:04:05Z") != now.Add(3*time.Second).Format("2006-01-02T15:04:05Z") {
		t.Errorf("%v is not equal %v", creds.RefreshExpiresAt, now.Add(3*time.Second).Format("2006-01-02T15:04:05Z"))
	}
}

type TokensAPIMock struct {
	tokensiface.TokensAPI
	GenerateResponse *tokens.GenerateResponse
	GenerateError    error
}

func (t *TokensAPIMock) Generate() (*tokens.GenerateResponse, error) {
	return t.GenerateResponse, t.GenerateError
}

// There is tested only no credentials.
// Other patterns are tested in onelogin/credentials package.
func TestRefresh(t *testing.T) {
	var v *credentials.Value
	now := time.Now()
	a := &TokensAPIMock{
		GenerateResponse: &tokens.GenerateResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			CreatedAt:    now.Format("2006-01-02T15:04:05Z"),
			ExpiresIn:    10,
			AccountID:    1234567,
			TokenType:    "bearer",
		},
	}
	c := Config{
		Endpoint:     "endpoint",
		ClientToken:  "client-token",
		ClientSecret: "client-secret",
		Credentials:  credentials.New(a, v),
	}
	if err := c.Refresh(); err != nil {
		t.Errorf("%#v", err)
	}
	creds := c.Credentials.Credentials
	if creds.AccessToken != "access-token" {
		t.Errorf("%v is not equal %v", creds.AccessToken, "access-token")
	}
	if creds.RefreshToken != "refresh-token" {
		t.Errorf("%v is not equal %v", creds.RefreshToken, "refresh-token")
	}
	if creds.CreatedAt.Format("2006-01-02T15:04:05Z") != now.Format("2006-01-02T15:04:05Z") {
		t.Errorf("%v is not equal %v", creds.CreatedAt, now.Format("2006-01-02T15:04:05Z"))
	}
	if creds.AccessExpiresAt.Format("2006-01-02T15:04:05Z") != now.Add(10*time.Second).Format("2006-01-02T15:04:05Z") {
		t.Errorf("%v is not equal %v", creds.AccessExpiresAt, now.Add(10*time.Second).Format("2006-01-02T15:04:05Z"))
	}
	if creds.RefreshExpiresAt.Format("2006-01-02T15:04:05Z") != now.Add(45*24*time.Hour).Format("2006-01-02T15:04:05Z") {
		t.Errorf("%v is not equal %v", creds.RefreshExpiresAt, now.Add(45*24*time.Hour).Format("2006-01-02T15:04:05Z"))
	}
}

func TestSave(t *testing.T) {
	CacheDir = os.TempDir()
	var file = path.Join(CacheDir, fmt.Sprintf("onelogin.%s.cache", "client-token"))
	defer os.Remove(file)
	var v *credentials.Value
	now := time.Now()
	a := &TokensAPIMock{
		GenerateResponse: &tokens.GenerateResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			CreatedAt:    now.Format("2006-01-02T15:04:05Z"),
			ExpiresIn:    10,
			AccountID:    1234567,
			TokenType:    "bearer",
		},
	}
	c := Config{
		Endpoint:     "endpoint",
		ClientToken:  "client-token",
		ClientSecret: "client-secret",
		Credentials:  credentials.New(a, v),
	}
	if err := c.Save(); err != nil {
		t.Errorf("%#v", err)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	actual := string(data)
	expected := fmt.Sprintf(`AccessToken = "access-token"
RefreshToken = "refresh-token"
CreatedAt = %s
AccessExpiresAt = %s
RefreshExpiresAt = %s
`,
		now.Format("2006-01-02T15:04:05Z"),
		now.Add(10*time.Second).Format("2006-01-02T15:04:05Z"),
		now.Add(45*24*time.Hour).Format("2006-01-02T15:04:05Z"),
	)
	if actual != expected {
		t.Errorf("%s is not equal %s", actual, expected)
	}
}

func TestSaveError(t *testing.T) {
	CacheDir = os.TempDir()
	var file = path.Join(CacheDir, fmt.Sprintf("onelogin.%s.cache", "client-token"))
	defer os.Remove(file)
	var v *credentials.Value
	a := &TokensAPIMock{
		GenerateError: fmt.Errorf("generate error"),
	}
	c := Config{
		Endpoint:     "endpoint",
		ClientToken:  "client-token",
		ClientSecret: "client-secret",
		Credentials:  credentials.New(a, v),
	}
	if err := c.Save(); err.Error() != "generate error" {
		t.Errorf("%#v", err)
	}
	info, err := os.Stat(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	if info.Size() > 0 {
		t.Error("file size is not zero")
	}
}
