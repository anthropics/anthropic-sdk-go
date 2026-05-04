package auth

import (
	"fmt"
	"net/url"
	"strings"
)

// requireSecureTokenEndpoint rejects base URLs that would cause a JWT
// assertion or refresh token to be sent over cleartext HTTP. Local
// development hosts (localhost, 127.0.0.1, ::1) are allowed.
func requireSecureTokenEndpoint(base string) error {
	if base == "" {
		return nil
	}
	u, err := url.Parse(base)
	if err != nil {
		return fmt.Errorf("invalid token endpoint base URL %q: %w", base, err)
	}
	if u.Scheme == "https" {
		return nil
	}
	if u.Scheme == "http" && isLoopbackHost(u.Hostname()) {
		return nil
	}
	return fmt.Errorf("refusing to send credential over non-https token endpoint %q", base)
}

func isLoopbackHost(host string) bool {
	host = strings.ToLower(host)
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
