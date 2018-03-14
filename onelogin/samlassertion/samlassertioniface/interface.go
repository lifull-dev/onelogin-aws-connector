package samlassertioniface

import "github.com/lifull-dev/onelogin-aws-connector/onelogin/samlassertion"

// SAMLAssertionAPI is SAMLAssertion API Interface
type SAMLAssertionAPI interface {
	Generate(input *samlassertion.GenerateRequest) (*samlassertion.GenerateResponse, error)
}
