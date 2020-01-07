package login

import (
	"strconv"

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

func (l *Login) Login(choice func(devices []samlassertion.GenerateResponseFactorDevice) (int, error), input func() (string, error)) (*sts.Credentials, error) {
	assertion, err := l.generateAssertion()
	if err != nil {
		return nil, err
	}
	SAML := assertion.SAML
	if SAML == "" {
		factor := assertion.Factors[0]
		selected := 0
		length := len(factor.Devices)
		if length > 1 {
			selected, err = choice(factor.Devices)
			if err != nil {
				return nil, err
			}
		}
		deviceID := factor.Devices[selected].DeviceID
		token, err := input()
		if err != nil {
			return nil, err
		}
		verified, err := l.generateAssertionWithMFA(deviceID, factor.StateToken, token)
		if err != nil {
			return nil, err
		}
		SAML = verified.SAML
	}
	return l.assumeRole(SAML)
}

// Execute represents login flow
func (l *Login) generateAssertion() (*samlassertion.GenerateResponse, error) {
	input := &samlassertion.GenerateRequest{
		UsernameOrEmail: l.Params.UsernameOrEmail,
		Password:        l.Params.Password,
		AppID:           l.Params.AppID,
		Subdomain:       l.Params.Subdomain,
	}
	return l.SAMLAssertion.Generate(input)
}

func (l *Login) generateAssertionWithMFA(deviceId int, stateToken string, otpToken string) (*samlassertion.VerifyFactorResponse, error) {
	input := &samlassertion.VerifyFactorRequest{
		AppID:       l.Params.AppID,
		DeviceID:    strconv.Itoa(deviceId),
		StateToken:  stateToken,
		OtpToken:    otpToken,
		DoNotNotify: false,
	}
	return l.SAMLAssertion.VerifyFactor(input)
}

// Execute represents login flow
func (l *Login) assumeRole(SAML string) (*sts.Credentials, error) {
	assumeRoleInput := &sts.AssumeRoleWithSAMLInput{
		PrincipalArn:    &l.Params.PrincipalArn,
		RoleArn:         &l.Params.RoleArn,
		SAMLAssertion:   &SAML,
		DurationSeconds: &l.Params.DurationSeconds,
	}
	assumeRoleOutput, err := l.STS.AssumeRoleWithSAML(assumeRoleInput)
	if err != nil {
		return nil, err
	}
	return assumeRoleOutput.Credentials, nil
}
