package login

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/lifull-dev/onelogin-aws-connector/onelogin"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/samlassertion"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/samlassertion/samlassertioniface"
)

// Login represents login
type Login struct {
	SAMLAssertion samlassertioniface.SAMLAssertionAPI
	STS           stsiface.STSAPI
	Params        *Parameters
}

// Parameters represents login parameters
type Parameters struct {
	UsernameOrEmail string
	Password        string
	AppID           string
	Subdomain       string
	PrincipalArn    string
	RoleArn         string
	DurationSeconds int64
}

// New creates a Login instance
func New(config *onelogin.Config, params *Parameters) *Login {
	return &Login{
		SAMLAssertion: samlassertion.NewSAMLAssertion(config),
		STS:           sts.New(session.New()),
		Params:        params,
	}
}

// Execute represents login flow
func (l *Login) Execute() (*sts.Credentials, error) {
	input := &samlassertion.GenerateRequest{
		UsernameOrEmail: l.Params.UsernameOrEmail,
		Password:        l.Params.Password,
		AppID:           l.Params.AppID,
		Subdomain:       l.Params.Subdomain,
	}
	saml, err := l.SAMLAssertion.Generate(input)
	if err != nil {
		return nil, err
	}

	assumeRoleInput := &sts.AssumeRoleWithSAMLInput{
		PrincipalArn:    &l.Params.PrincipalArn,
		RoleArn:         &l.Params.RoleArn,
		SAMLAssertion:   &saml.SAML,
		DurationSeconds: &l.Params.DurationSeconds,
	}
	assumeRoleOutput, err := l.STS.AssumeRoleWithSAML(assumeRoleInput)
	if err != nil {
		return nil, err
	}
	return assumeRoleOutput.Credentials, nil
}
