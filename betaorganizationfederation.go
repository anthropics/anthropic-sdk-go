// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/option"
)

// BetaOrganizationFederationService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationFederationService] method instead.
type BetaOrganizationFederationService struct {
	Options []option.RequestOption
	Issuers BetaOrganizationFederationIssuerService
	Rules   BetaOrganizationFederationRuleService
}

// NewBetaOrganizationFederationService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationFederationService(opts ...option.RequestOption) (r BetaOrganizationFederationService) {
	r = BetaOrganizationFederationService{}
	r.Options = opts
	r.Issuers = NewBetaOrganizationFederationIssuerService(opts...)
	r.Rules = NewBetaOrganizationFederationRuleService(opts...)
	return
}
