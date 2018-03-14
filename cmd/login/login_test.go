package login

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin/samlassertion"
)

type SAMLAssertionMock struct {
	GenerateResponse *samlassertion.GenerateResponse
	Error            error
	InputVerifier    func(*samlassertion.GenerateRequest) error
}

func (s *SAMLAssertionMock) Generate(*samlassertion.GenerateRequest) (*samlassertion.GenerateResponse, error) {
	return s.GenerateResponse, s.Error
}

type STSMock struct {
	stsiface.STSAPI
	AssumeRoleWithSAMLOutput *sts.AssumeRoleWithSAMLOutput
	Error                    error
	InputVerifier            func(*sts.AssumeRoleWithSAMLInput) error
}

func (s *STSMock) AssumeRoleWithSAML(*sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
	return s.AssumeRoleWithSAMLOutput, s.Error
}

func TestExecuteSuccess(t *testing.T) {
	now := time.Now()
	l := &Login{
		SAMLAssertion: &SAMLAssertionMock{
			GenerateResponse: &samlassertion.GenerateResponse{
				Status: &samlassertion.GenerateResponseStatus{},
				SAML:   "base64-encoded-saml-data",
			},
			InputVerifier: func(request *samlassertion.GenerateRequest) error {
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
		},
		STS: &STSMock{
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
				if *request.SAMLAssertion != "base64-encoded-saml-data" {
					t.Errorf("%s is not equal %s", *request.SAMLAssertion, "base64-encoded-saml-data")
				}
				return nil
			},
		},
		Params: &Parameters{
			UsernameOrEmail: "username-or-email",
			Password:        "password",
			AppID:           "app-id",
			Subdomain:       "subdomain",
			PrincipalArn:    "principal-arn",
			RoleArn:         "role-arn",
		},
	}
	creds, err := l.Execute()
	if err != nil {
		t.Errorf("%v", err)
	}
	if *creds.AccessKeyId != "access-key-id" {
		t.Errorf("%s is not equal %s", *creds.AccessKeyId, "access-key-id")
	}
	if *creds.SecretAccessKey != "secret-access-key" {
		t.Errorf("%s is not equal %s", *creds.SecretAccessKey, "secret-access-key")
	}
	if *creds.SessionToken != "session-token" {
		t.Errorf("%s is not equal %s", *creds.SessionToken, "session-token")
	}
	if creds.Expiration.String() != now.String() {
		t.Errorf("%s is not equal %s", creds.Expiration.String(), now.String())
	}
}

func TestExecuteErrorOnSAMLAssertion(t *testing.T) {
	l := &Login{
		SAMLAssertion: &SAMLAssertionMock{
			Error: fmt.Errorf("SAMLAssertion Error"),
			InputVerifier: func(request *samlassertion.GenerateRequest) error {
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
		},
		STS: &STSMock{},
		Params: &Parameters{
			UsernameOrEmail: "username-or-email",
			Password:        "password",
			AppID:           "app-id",
			Subdomain:       "subdomain",
			PrincipalArn:    "principal-arn",
			RoleArn:         "role-arn",
		},
	}
	_, err := l.Execute()
	if err.Error() != "SAMLAssertion Error" {
		t.Errorf("%v", err)
	}
}

func TestExecuteErrorOnAssumeRole(t *testing.T) {
	l := &Login{
		SAMLAssertion: &SAMLAssertionMock{
			GenerateResponse: &samlassertion.GenerateResponse{
				Status: &samlassertion.GenerateResponseStatus{},
				SAML:   "base64-encoded-saml-data",
			},
			InputVerifier: func(request *samlassertion.GenerateRequest) error {
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
		},
		STS: &STSMock{
			Error: fmt.Errorf("AssumeRole Error"),
			InputVerifier: func(request *sts.AssumeRoleWithSAMLInput) error {
				if *request.PrincipalArn != "principal-arn" {
					t.Errorf("%s is not equal %s", *request.PrincipalArn, "principal-arn")
				}
				if *request.RoleArn != "role-arn" {
					t.Errorf("%s is not equal %s", *request.RoleArn, "role-arn")
				}
				if *request.SAMLAssertion != "base64-encoded-saml-data" {
					t.Errorf("%s is not equal %s", *request.SAMLAssertion, "base64-encoded-saml-data")
				}
				return nil
			},
		},
		Params: &Parameters{
			UsernameOrEmail: "username-or-email",
			Password:        "password",
			AppID:           "app-id",
			Subdomain:       "subdomain",
			PrincipalArn:    "principal-arn",
			RoleArn:         "role-arn",
		},
	}
	_, err := l.Execute()
	if err.Error() != "AssumeRole Error" {
		t.Errorf("%v", err)
	}
}

func StringRef(v string) *string {
	return &v
}
