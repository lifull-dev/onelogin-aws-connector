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
			if request.DoNotNotify != false {
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
			t.Errorf("%s is not equal %s", request.OtpToken, "123456")
		}
		if request.DoNotNotify != false {
			t.Errorf("%v is not equal %v", request.DoNotNotify, false)
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
	_, err := l.Login(func(devices []samlassertion.GenerateResponseFactorDevice) (i int, err error) {
		t.Errorf("%s", "Don't call choice function")
		return 0, nil
	}, func() (s string, err error) {
		t.Errorf("%s", "Don't call mfa function")
		return "", nil
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
	_, err := l.Login(func(devices []samlassertion.GenerateResponseFactorDevice) (i int, err error) {
		t.Errorf("%s", "Don't call choice function")
		return 0, nil
	}, func() (s string, err error) {
		return "", nil
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
	_, err := l.Login(func(devices []samlassertion.GenerateResponseFactorDevice) (i int, err error) {
		t.Errorf("%s", "Don't call choice function")
		return 0, nil
	}, func() (s string, err error) {
		return "765432", nil
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
	_, err := l.Login(func(devices []samlassertion.GenerateResponseFactorDevice) (i int, err error) {
		return 1, nil
	}, func() (s string, err error) {
		return "098765", nil
	})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestLogin_LoginChoiceErrorWithMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionForMultipleMFA(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(func(devices []samlassertion.GenerateResponseFactorDevice) (i int, err error) {
		return 0, errors.Errorf("choice error")
	}, func() (s string, err error) {
		return "", nil
	})
	if err != nil && err.Error() != "choice error" {
		t.Errorf("'%s' is not equal 'choice error'", err.Error())
	}
}

func TestLogin_LoginMFAErrorWithMFA(t *testing.T) {
	l := &Login{
		SAMLAssertion: createAssertionForSingleMFA(t),
		STS:           createSTS(t),
		Params:        createDefaultParams(),
	}
	_, err := l.Login(func(devices []samlassertion.GenerateResponseFactorDevice) (i int, err error) {
		t.Errorf("%s", "Don't call choice function")
		return 0, nil
	}, func() (s string, err error) {
		return "", errors.Errorf("mfa error")
	})
	if err != nil && err.Error() != "mfa error" {
		t.Errorf("'%s' is not equal 'mfa error'", err.Error())
	}
}

func StringRef(v string) *string {
	return &v
}
