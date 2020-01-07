package samlassertion

import (
	"encoding/json"

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

// GenerateTemporaryResponse response of OneLogin Generate Tokens v2 API
type GenerateResponse struct {
	Status  *GenerateResponseStatus `json:"status"`
	SAML    string
	Factors []GenerateResponseFactor
}

// GenerateResponse response of OneLogin Generate Tokens v2 API
type GenerateSAMLResponse struct {
	Status *GenerateResponseStatus `json:"status"`
	SAML   string                  `json:"data"`
}

// GenerateFactorsResponse response of OneLogin Generate Tokens v2 API
type GenerateFactorsResponse struct {
	Status  *GenerateResponseStatus  `json:"status"`
	Factors []GenerateResponseFactor `json:"data"`
}

// GenerateResponseStatus status
type GenerateResponseStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
}

type GenerateResponseFactor struct {
	StateToken  string                         `json:"state_token"`
	Devices     []GenerateResponseFactorDevice `json:"devices"`
	CallbackURL string                         `json:"callback_url"`
	User        *GenerateResponseFactorUser    `json:"user"`
}

type GenerateResponseFactorDevice struct {
	DeviceID   int    `json:"device_id"`
	DeviceType string `json:"device_type"`
}

type GenerateResponseFactorUser struct {
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	ID        int    `json:"id"`
}

// Generate call generate tokens v2
func (s *SAMLAssertion) Generate(input *GenerateRequest) (*GenerateResponse, error) {
	inputJSON, err := json.Marshal(input)
	body, err := s.post("/api/1/saml_assertion", inputJSON)
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
	if output.Status.Message == "Success" {
		var saml GenerateSAMLResponse
		if err := json.Unmarshal(body, &saml); err != nil {
			return nil, err
		}
		output.SAML = saml.SAML
	} else {
		var factors GenerateFactorsResponse
		if err := json.Unmarshal(body, &factors); err != nil {
			return nil, err
		}
		output.Factors = factors.Factors
	}
	return &output, nil
}
