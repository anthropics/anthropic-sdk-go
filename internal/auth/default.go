package auth

import (
	"os"
	"strings"
)

// identityProviderFromEnv resolves an [IdentityTokenProvider] from env vars,
// preferring ANTHROPIC_IDENTITY_TOKEN_FILE (rotation-safe) over the literal
// ANTHROPIC_IDENTITY_TOKEN. Returns nil when neither is set.
func identityProviderFromEnv() IdentityTokenProvider {
	if path, ok := os.LookupEnv(EnvIdentityTokenFile); ok {
		return &IdentityTokenFile{Path: path}
	}
	if token, ok := os.LookupEnv(EnvIdentityToken); ok {
		return &IdentityTokenValue{Token: token}
	}
	return nil
}

// EnvCredentials resolves OIDC-federation credentials from environment
// variables. Returns a usable result when the three required vars are all
// set. Otherwise returns nil plus a detail + state enum: NotSet when
// nothing is configured, Partial with a "missing: ..." detail otherwise.
func EnvCredentials() (*CredentialsResult, string, CredentialSourceState) {
	federationRuleID, hasFedRule := os.LookupEnv(EnvFederationRuleID)
	_, hasOrgID := os.LookupEnv(EnvOrganizationID)
	hasTokenFile := os.Getenv(EnvIdentityTokenFile) != ""
	hasTokenLiteral := os.Getenv(EnvIdentityToken) != ""
	hasToken := hasTokenFile || hasTokenLiteral

	if !hasFedRule && !hasOrgID && !hasToken {
		return nil, "", CredentialSourceNotSet
	}

	if !hasFedRule || !hasOrgID || !hasToken {
		var missing []string
		if !hasFedRule {
			missing = append(missing, EnvFederationRuleID)
		}
		if !hasOrgID {
			missing = append(missing, EnvOrganizationID)
		}
		if !hasToken {
			missing = append(missing, EnvIdentityTokenFile+" (or "+EnvIdentityToken+")")
		}
		return nil, "partial configuration, missing: " + strings.Join(missing, ", "), CredentialSourcePartial
	}

	orgID := os.Getenv(EnvOrganizationID)
	identityProvider := identityProviderFromEnv()
	cfg := OIDCFederationConfig{
		IdentityProvider: identityProvider,
		FederationRuleID: federationRuleID,
		OrganizationID:   orgID,
	}
	if sa, ok := os.LookupEnv(EnvServiceAccountID); ok {
		cfg.ServiceAccountID = sa
	}
	if ws, ok := os.LookupEnv(EnvWorkspaceID); ok {
		// An empty ANTHROPIC_WORKSPACE_ID (a defaulted-but-empty CI variable)
		// is treated as unset: WorkspaceID stays "" and the wire field has
		// json:"workspace_id,omitempty", so `"workspace_id": ""` is never
		// serialized. No coercion needed — Go's Getenv + omitempty handle it.
		cfg.WorkspaceID = ws
	}

	return &CredentialsResult{
		Provider: NewOIDCFederationCredentials(cfg),
	}, "", CredentialSourceNotSet
}
