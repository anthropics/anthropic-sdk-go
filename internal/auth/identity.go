package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// IdentityTokenFile reads a JWT from a file on each call, supporting
// automatic rotation (e.g. Kubernetes projected service account tokens).
type IdentityTokenFile struct {
	Path string
}

func (p *IdentityTokenFile) GetIdentityToken(_ context.Context) (string, error) {
	if p.Path == "" {
		return "", fmt.Errorf("identity token file path is empty")
	}
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read identity token file %q: %w", p.Path, err)
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", fmt.Errorf("identity token file %q is empty", p.Path)
	}
	return token, nil
}

// IdentityTokenValue wraps a static JWT string as an [IdentityTokenProvider].
type IdentityTokenValue struct {
	Token string
}

func (p *IdentityTokenValue) GetIdentityToken(_ context.Context) (string, error) {
	if p.Token == "" {
		return "", fmt.Errorf("identity token value is empty")
	}
	return p.Token, nil
}
