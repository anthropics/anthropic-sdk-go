package auth

import (
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
)

// OAuth wire-contract constants live in the public config package so the
// SDK has a single source of truth. Re-exported here as local aliases so
// auth code can reference them without a package qualifier.
const (
	GrantTypeJWTBearer    = config.GrantTypeJWTBearer
	GrantTypeRefreshToken = config.GrantTypeRefreshToken

	TokenEndpoint = config.TokenEndpoint

	// FederationBetaHeader is sent on requests to the oauth/token endpoint
	// during federation token exchange.
	FederationBetaHeader = config.FederationBetaHeader

	// OAuthAPIBetaHeader is sent on authenticated API requests that use a
	// bearer token obtained via OAuth/federation, and on refresh_token
	// grants against the token endpoint.
	OAuthAPIBetaHeader = config.OAuthAPIBetaHeader
)

const (
	AdvisoryRefreshThreshold  = 120 * time.Second
	MandatoryRefreshThreshold = 30 * time.Second

	EnvIdentityToken     = "ANTHROPIC_IDENTITY_TOKEN"
	EnvIdentityTokenFile = "ANTHROPIC_IDENTITY_TOKEN_FILE"
	EnvFederationRuleID  = "ANTHROPIC_FEDERATION_RULE_ID"
	EnvOrganizationID    = "ANTHROPIC_ORGANIZATION_ID"
	EnvServiceAccountID  = "ANTHROPIC_SERVICE_ACCOUNT_ID"
	EnvWorkspaceID       = "ANTHROPIC_WORKSPACE_ID"
)
