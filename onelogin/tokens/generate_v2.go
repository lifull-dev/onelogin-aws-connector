package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// https://developers.onelogin.com/api-docs/1/oauth20-tokens/generate-tokens-2

// GenerateRequest request for OneLogin Generate Tokens v2 API
type GenerateRequest struct {
	GrantType string `json:"grant_type"`
}

// RefreshRequest request for OneLogin Generate Tokens v2 API
type RefreshRequest struct {
	GrantType    string `json:"grant_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateResponse response of OneLogin Generate Tokens v2 API
type GenerateResponse struct {
	Status       *Status `json:"status"`
	AccessToken  string  `json:"access_token"`
	CreatedAt    string  `json:"created_at"`
	ExpiresIn    int     `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	AccountID    int     `json:"account_id"`
}

// RefreshResponse response of OneLogin Refresh Tokens v2 API
type RefreshResponse = GenerateResponse

// Status status
type Status struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
}

// Tokens OneLogin Generate Tokens v2 API
type Tokens struct {
	Endpoint     string
	ClientToken  string
	ClientSecret string
	HTTPClient   *http.Client
}

// NewTokens creates a Tokens
func NewTokens() *Tokens {
	return &Tokens{
		HTTPClient: &http.Client{},
	}
}

// Generate retrive access_token and other
func (g *Tokens) Generate() (*GenerateResponse, error) {
	input := &GenerateRequest{
		GrantType: "client_credentials",
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s/auth/oauth2/v2/token", g.Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(inputJSON)))
	if err != nil {
		return nil, err
	}
	creds := fmt.Sprintf("client_id:%s, client_secret:%s", g.ClientToken, g.ClientSecret)
	req.Header.Set("Authorization", creds)
	req.Header.Set("Content-Type", "application/json")
	client := g.HTTPClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var output GenerateResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}
	if output.Status != nil && output.Status.Error {
		return nil, errors.Errorf("(%d) %s", output.Status.Code, output.Status.Message)
	}
	return &output, nil
}

// Refresh retrive access_token and other by refresh_token
func (g *Tokens) Refresh(input *RefreshRequest) (*RefreshResponse, error) {
	input.GrantType = "refresh_token"
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s/auth/oauth2/v2/token", g.Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(inputJSON)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := g.HTTPClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var output RefreshResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}
	if output.Status != nil && output.Status.Error {
		return nil, errors.Errorf("[%d] %s: %s", output.Status.Code, output.Status.Type, output.Status.Message)
	}
	return &output, nil
}
