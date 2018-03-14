package samlassertion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin"
	"github.com/pkg/errors"
)

// https://developers.onelogin.com/api-docs/1/saml-assertions/generate-saml-assertion

// GenerateRequest request for OneLogin Generate Tokens v2 API
type GenerateRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
	AppID           string `json:"app_id"`
	Subdomain       string `json:"subdomain"`
	IPAddress       string `json:"ip_address"`
}

// GenerateResponse response of OneLogin Generate Tokens v2 API
type GenerateResponse struct {
	Status *GenerateResponseStatus `json:"status"`
	SAML   string                  `json:"data"`
}

// GenerateResponseStatus status
type GenerateResponseStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
}

// SAMLAssertion OneLogin Generate SAML Assertion API
type SAMLAssertion struct {
	config     *onelogin.Config
	HTTPClient *http.Client
}

// NewSAMLAssertion creates a SAMLAssertion
func NewSAMLAssertion(config *onelogin.Config) *SAMLAssertion {
	return &SAMLAssertion{
		config:     config,
		HTTPClient: &http.Client{},
	}
}

// Generate call generate tokens v2
func (s *SAMLAssertion) Generate(input *GenerateRequest) (*GenerateResponse, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s/api/1/saml_assertion", s.config.Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(inputJSON)))
	if err != nil {
		return nil, err
	}
	creds, err := s.config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	authorization := fmt.Sprintf("bearer:%s", creds.AccessToken)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
	client := s.HTTPClient
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
	if output.Status.Error {
		return nil, errors.Errorf("[%d] %s: %s", output.Status.Code, output.Status.Type, output.Status.Message)
	}
	return &output, nil
}
