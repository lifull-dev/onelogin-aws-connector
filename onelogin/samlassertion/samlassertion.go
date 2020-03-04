package samlassertion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin"
)

// SAMLAssertion OneLogin Generate SAML Assertion API
type SAMLAssertion struct {
	config                   *onelogin.Config
	HTTPClient               *http.Client
	verifyFactorLoopMax      int
	verifyFactorLoopDuration int
}

// https://developers.onelogin.com/api-docs/1/saml-assertions/generate-saml-assertion

// GenerateRequest request for OneLogin Generate Tokens v2 API
type GenerateRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
	AppID           string `json:"app_id"`
	Subdomain       string `json:"subdomain"`
	IPAddress       string `json:"ip_address"`
}

// GenerateResponse response
type GenerateResponse struct {
	Status  *GenerateResponseStatus `json:"status"`
	SAML    string
	Factors []GenerateResponseFactor
}

// GenerateSAMLResponse response of OneLogin Generate Tokens v2 API without mfa
type GenerateSAMLResponse struct {
	Status *GenerateResponseStatus `json:"status"`
	SAML   string                  `json:"data"`
}

// GenerateFactorsResponse response of OneLogin Generate Tokens v2 API with mfa
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
	DeviceID        int    `json:"device_id"`
	DeviceType      string `json:"device_type"`
	RequireOTPToken bool
}

type GenerateResponseFactorUser struct {
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	ID        int    `json:"id"`
}

// https://developers.onelogin.com/api-docs/1/saml-assertions/verify-factor

// VerifyFactorRequest request for OneLogin VerifyFactor Tokens v2 API
type VerifyFactorRequest struct {
	AppID       string `json:"app_id"`
	DeviceID    string `json:"device_id"`
	StateToken  string `json:"state_token"`
	OtpToken    string `json:"otp_token"`
	DoNotNotify bool   `json:"do_not_notify"`
}

// VerifyFactorTemporaryResponse response of OneLogin VerifyFactor Tokens v2 API
type VerifyFactorResponse struct {
	Status *VerifyFactorResponseStatus `json:"status"`
	SAML   string                      `json:"data"`
}

// VerifyFactorResponseStatus status
type VerifyFactorResponseStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
}

// NewSAMLAssertion creates a SAMLAssertion
func NewSAMLAssertion(config *onelogin.Config) *SAMLAssertion {
	return &SAMLAssertion{
		config:                   config,
		HTTPClient:               &http.Client{},
		verifyFactorLoopMax:      60,
		verifyFactorLoopDuration: 1000,
	}
}

// Generate call generate tokens v2
func (s *SAMLAssertion) Generate(input *GenerateRequest) (*GenerateResponse, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
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
		devices := factors.Factors[0].Devices
		for i := range devices {
			devices[i].RequireOTPToken = true
			device := devices[i]
			if device.DeviceType == "OneLogin Protect" {
				devices = append(devices, GenerateResponseFactorDevice{
					DeviceType:      "Notify to OneLogin Protect",
					DeviceID:        device.DeviceID,
					RequireOTPToken: false,
				})
			}
		}
		factors.Factors[0].Devices = devices
		output.Factors = factors.Factors
	}
	return &output, nil
}

// VerifyFactor call VerifyFactor tokens v2
func (s *SAMLAssertion) VerifyFactor(input *VerifyFactorRequest) (*VerifyFactorResponse, error) {
	return s.verifyFactor(input, 0)
}

func (s *SAMLAssertion) verifyFactor(input *VerifyFactorRequest, loopCount int) (*VerifyFactorResponse, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	body, err := s.post("/api/1/saml_assertion/verify_factor", inputJSON)
	if err != nil {
		return nil, err
	}
	var output VerifyFactorResponse
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}
	if output.Status.Error {
		return nil, errors.Errorf("[%d] %s: %s", output.Status.Code, output.Status.Type, output.Status.Message)
	}
	if output.Status.Type == "pending" {
		if loopCount >= s.verifyFactorLoopMax {
			return nil, errors.Errorf("[%d] timed out: %s", output.Status.Code, output.Status.Message)
		}
		time.Sleep(time.Duration(s.verifyFactorLoopDuration))
		input.DoNotNotify = true
		return s.verifyFactor(input, loopCount+1)
	}
	return &output, nil
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
