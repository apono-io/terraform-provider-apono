package client

import (
	"context"
)

// TokenSecuritySource implements SecuritySource with a static token.
type TokenSecuritySource struct {
	token string
}

var _ SecuritySource = &TokenSecuritySource{}

func NewTokenSecuritySource(token string) *TokenSecuritySource {
	return &TokenSecuritySource{
		token: token,
	}
}

func (s *TokenSecuritySource) Authorization(_ context.Context, _ OperationName) (Authorization, error) {
	return Authorization{
		Token: s.token,
	}, nil
}
