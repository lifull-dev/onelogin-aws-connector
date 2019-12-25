package samlassertion

import (
	"encoding/json"
	"github.com/pkg/errors"
)

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

// VerifyFactor call VerifyFactor tokens v2
func (s *SAMLAssertion) VerifyFactor(input *VerifyFactorRequest) (*VerifyFactorResponse, error) {
	inputJSON, err := json.Marshal(input)
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
	return &output, nil
}
