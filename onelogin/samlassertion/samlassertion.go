package samlassertion

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin"
)

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

// post OneLogin API Request
func (s *SAMLAssertion) post(path string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("https://%s%s", s.config.Endpoint, path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	credentials, err := s.config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	authorization := fmt.Sprintf("bearer:%s", credentials.AccessToken)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
	client := s.HTTPClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
