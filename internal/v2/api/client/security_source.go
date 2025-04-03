package client

import (
	"context"
)

// TokenSecuritySource implements SecuritySource with a static token.
type TokenSecuritySource struct {
	token string
}

// NewTokenSecuritySource creates a new TokenSecuritySource with the given token.
func NewTokenSecuritySource(token string) *TokenSecuritySource {
	return &TokenSecuritySource{
		token: token,
	}
}

// Authorization implements SecuritySource.Authorization.
func (s *TokenSecuritySource) Authorization(_ context.Context, _ OperationName) (Authorization, error) {
	return Authorization{
		Token: s.token,
	}, nil
}
