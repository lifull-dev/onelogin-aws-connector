package tokensiface

import "github.com/lifull-dev/onelogin-aws-connector/onelogin/tokens"

// TokensAPI is Tokens API Interface
type TokensAPI interface {
	Generate() (*tokens.GenerateResponse, error)
	Refresh(input *tokens.RefreshRequest) (*tokens.RefreshResponse, error)
}
