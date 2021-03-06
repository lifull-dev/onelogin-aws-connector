package login

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin/samlassertion"
)

type SAMLAssertionMock struct {
	GenerateResponse          *samlassertion.GenerateResponse
	GenerateInputVerifier     func(*samlassertion.GenerateRequest) error
	GenerateError             error
	VerifyFactorResponse      *samlassertion.VerifyFactorResponse
	VerifyFactorInputVerifier func(request *samlassertion.VerifyFactorRequest) error
	VerifyFactorError         error
}

func (s *SAMLAssertionMock) Generate(request *samlassertion.GenerateRequest) (*samlassertion.GenerateResponse, error) {
	if err := s.GenerateInputVerifier(request); err != nil {
		return nil, err
	}
	return s.GenerateResponse, s.GenerateError
}

func (s *SAMLAssertionMock) VerifyFactor(request *samlassertion.VerifyFactorRequest) (*samlassertion.VerifyFactorResponse, error) {
	if err := s.VerifyFactorInputVerifier(request); err != nil {
		return nil, err
	}
	return s.VerifyFactorResponse, s.VerifyFactorError
}

type STSMock struct {
	stsiface.STSAPI
	AssumeRoleWithSAMLOutput *sts.AssumeRoleWithSAMLOutput
	Error                    error
	InputVerifier            func(*sts.AssumeRoleWithSAMLInput) error
}

func (s *STSMock) AssumeRoleWithSAML(input *sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
	if err := s.InputVerifier(input); err != nil {
		return nil, err
	}
	return s.AssumeRoleWithSAMLOutput, s.Error
}

type EventMock struct {
	DeviceIndex int
	ChooseError error
	MFAToken    string
	InputError  error
}

func (m *EventMock) ChooseDeviceIndex(devices []samlassertion.GenerateResponseFactorDevice) (int, error) {
	return m.DeviceIndex, m.ChooseError
}
func (m *EventMock) InputMFAToken() (string, error) {
	return m.MFAToken, m.InputError
}

func createAssertion(t *testing.T) *SAMLAssertionMock {
	return &SAMLAssertionMock{
		GenerateResponse: &samlassertion.GenerateResponse{
			SAML: "Base64 encoded SAML Data",
		},
		GenerateInputVerifier: func(request *samlassertion.GenerateRequest) error {
			if request.UsernameOrEmail != "username-or-email" {
				t.Errorf("%s is not equal %s", request.UsernameOrEmail, "username-or-email")
			}
			if request.Password != "password" {
				t.Errorf("%s is not equal %s", request.Password, "password")
			}
			if request.AppID != "app-id" {
				t.Errorf("%s is not equal %s", request.AppID, "app-id")
			}
			if request.Subdomain != "subdomain" {
				t.Errorf("%s is not equal %s", request.Subdomain, "subdomain")
			}
			return nil
		},
	}
}

func createAssertionForSingleMFA(t *testing.T) *SAMLAssertionMock {
	return &SAMLAssertionMock{
		GenerateResponse: &samlassertion.GenerateResponse{
			Factors: []samlassertion.GenerateResponseFactor{
				{
					StateToken: "state-token",
					Devices: []samlassertion.GenerateResponseFactorDevice{
						{
							DeviceID:   345678,
							DeviceType: "device type 1",
							RequireOTPToken: true,
						},
					},
				},
			},
		},
		GenerateInputVerifier: func(request *samlassertion.GenerateRequest) error {
			return nil
		},
		VerifyFactorResponse: &samlassertion.VerifyFactorResponse{
			SAML: "Base64 encoded SAML Data",
		},
		VerifyFactorInputVerifier: func(request *samlassertion.VerifyFactorRequest) error {
			if request.AppID != "app-id" {
				t.Errorf("%s is not equal %s", request.AppID, "app-id")
			}
			if request.DeviceID != "345678" {
				t.Errorf("%s is not equal %s", request.DeviceID, "0")
			}
			if request.StateToken != "state-token" {
				t.Errorf("%s is not equal %s", request.StateToken, "state-token")
			}
			if request.OtpToken != "765432" {
				t.Errorf("%s is not equal %s", request.OtpToken, "123456")
			}
			if !request.DoNotNotify {
				t.Errorf("%v is not equal %v", request.DoNotNotify, false)
			}
			return nil
		},
	}
}

func createAssertionForMultipleMFA(t *testing.T) *SAMLAssertionMock {
	assertion := createAssertionForSingleMFA(t)
	assertion.GenerateResponse.Factors[0].Devices = append(
		assertion.GenerateResponse.Factors[0].Devices,
		samlassertion.GenerateResponseFactorDevice{
			DeviceID:   987654,
			DeviceType: "device type 2",
			RequireOTPToken: true,
		})
	assertion.VerifyFactorInputVerifier = func(request *samlassertion.VerifyFactorRequest) error {
		if request.AppID != "app-id" {
			t.Errorf("%s is not equal %s", request.AppID, "app-id")
		}
		if request.DeviceID != "987654" {
			t.Errorf("%s is not equal %s", request.DeviceID, "0")
		}
		if request.StateToken != "state-token" {
			t.Errorf("%s is not equal %s", request.StateToken, "state-token")
		}
		if request.OtpToken != "098765" {
			t.Errorf("%s is not equal %s", request.OtpToken, "098765")
		}
		if !request.DoNotNotify {
			t.Errorf("%v is not equal %v", request.DoNotNotify, false)
		}
		return nil
	}
	return assertion
}

func createAssertionForNotify(t *testing.T) *SAMLAssertionMock {
	assertion := createAssertionForSingleMFA(t)
	assertion.GenerateResponse.Factors[0].Devices = append(
		assertion.GenerateResponse.Factors[0].Devices,
		samlassertion.GenerateResponseFactorDevice{
			DeviceID:   987654,
			DeviceType: "Notify OneLogin Protect",
			RequireOTPToken: false,
		})
	assertion.VerifyFactorInputVerifier = func(request *samlassertion.VerifyFactorRequest) error {
		if request.AppID != "app-id" {
			t.Errorf("%s is not equal %s", request.AppID, "app-id")
		}
		if request.DeviceID != "987654" {
			t.Errorf("%s is not equal %s", request.DeviceID, "0")
		}
		if request.StateToken != "state-token" {
			t.Errorf("%s is not equal %s", request.StateToken, "state-token")
		}
		if request.OtpToken != "" {
			t.Errorf("'%s' is not equal '%s'", request.OtpToken, "")
		}
		if request.DoNotNotify {
			t.Errorf("%v is not equal %v", request.DoNotNotify, true)
		}
		return nil
	}
	return assertion
}

func createAssertionError(t *testing.T) *SAMLAssertionMock {
	return &SAMLAssertionMock{
		GenerateResponse: &samlassertion.GenerateResponse{},
		GenerateError:    fmt.Errorf("%s", "SAML Assertion Generate Error"),
		GenerateInputVerifier: func(request *samlassertion.GenerateRequest) error {
			if request.UsernameOrEmail != "username-or-email" {
				t.Errorf("%s is not equal %s", request.UsernameOrEmail, "username-or-email")
			}
			if request.Password != "password" {
				t.Errorf("%s is not equal %s", request.Password, "password")
			}
			if request.AppID != "app-id" {
				t.Errorf("%s is not equal %s", request.AppID, "app-id")
			}
			if request.Subdomain != "subdomain" {
				t.Errorf("%s is not equal %s", request.Subdomain, "subdomain")
			}
			return nil
		},
	}
}

func createSTS(t *testing.T) *STSMock {
	now := time.Now()
	return &STSMock{
		AssumeRoleWithSAMLOutput: &sts.AssumeRoleWithSAMLOutput{
			Credentials: &sts.Credentials{
				AccessKeyId:     StringRef("access-key-id"),
				SecretAccessKey: StringRef("secret-access-key"),
				SessionToken:    StringRef("session-token"),
				Expiration:      &now,
			},
		},
		InputVerifier: func(request *sts.AssumeRoleWithSAMLInput) error {
			if *request.PrincipalArn != "principal-arn" {
				t.Errorf("%s is not equal %s", *request.PrincipalArn, "principal-arn")
			}
			if *request.RoleArn != "role-arn" {
				t.Errorf("%s is not equal %s", *request.RoleArn, "role-arn")
			}
			if *request.SAMLAssertion != "Base64 encoded SAML Data" {
				t.Errorf("%s is not equal %s", *request.SAMLAssertion, "base64-encoded-saml-data")
			}
			return nil
		},
	}
}

func createDefaultParams() *Parameters {
	return &Parameters{
		UsernameOrEmail: "username-or-email",
		Password:        "password",
		AppID:           "app-id",
		Subdomain:       "subdomain",
		PrincipalArn:    "principal-arn",
		RoleArn:         "role-arn",
	}
}

func TestLogin_LoginWithoutMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertion(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(&EventMock{
		ChooseError: errors.New("Don't call choose function"),
		InputError:  errors.New("Don't call input function"),
	})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestLogin_LoginErrorWithoutMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionError(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(&EventMock{
		ChooseError: errors.New("Don't call choose function"),
	})
	if err != nil && err.Error() != "SAML Assertion Generate Error" {
		t.Errorf("'%s' is not equal 'SAML Assertion Generate Error'", err.Error())
	}
}

func TestLogin_LoginWithSingleMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionForSingleMFA(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(&EventMock{
		ChooseError: errors.New("Don't call choose function"),
		MFAToken:    "765432",
	})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestLogin_LoginWithMultipleMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionForMultipleMFA(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(&EventMock{
		DeviceIndex: 1,
		MFAToken:    "098765",
	})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestLogin_LoginWithNotify(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionForNotify(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(&EventMock{
		DeviceIndex: 1,
	})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestLogin_LoginChooseErrorWithMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionForMultipleMFA(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(&EventMock{
		ChooseError: errors.New("choose error"),
	})
	if err != nil && err.Error() != "choose error" {
		t.Errorf("'%s' is not equal 'choose error'", err.Error())
	}
}

func TestLogin_LoginMFAErrorWithMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionForSingleMFA(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(&EventMock{
		ChooseError: errors.New("Don't call choose function"),
		InputError:  errors.New("mfa error"),
	})
	if err != nil && err.Error() != "mfa error" {
		t.Errorf("'%s' is not equal 'mfa error'", err.Error())
	}
}

func StringRef(v string) *string {
	return &v
}
