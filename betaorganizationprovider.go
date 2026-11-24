package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaOrganizationProviderService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationProviderService] method instead.
type BetaOrganizationProviderService struct {
	Options []option.RequestOption
	Google  BetaOrganizationProviderGoogleService
}

// NewBetaOrganizationProviderService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationProviderService(opts ...option.RequestOption) (r BetaOrganizationProviderService) {
	r = BetaOrganizationProviderService{}
	r.Options = opts
	r.Google = NewBetaOrganizationProviderGoogleService(opts...)
	return
}
