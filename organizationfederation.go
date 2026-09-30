package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/option"
)

// OrganizationFederationService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationFederationService] method instead.
type OrganizationFederationService struct {
	Options []option.RequestOption
	Issuers OrganizationFederationIssuerService
	Rules   OrganizationFederationRuleService
}

// NewOrganizationFederationService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationFederationService(opts ...option.RequestOption) (r OrganizationFederationService) {
	r = OrganizationFederationService{}
	r.Options = opts
	r.Issuers = NewOrganizationFederationIssuerService(opts...)
	r.Rules = NewOrganizationFederationRuleService(opts...)
	return
}
