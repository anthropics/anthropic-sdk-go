package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaOrganizationProviderGoogleProjectService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationProviderGoogleProjectService] method instead.
type BetaOrganizationProviderGoogleProjectService struct {
	Options      []option.RequestOption
	DataControls BetaOrganizationProviderGoogleProjectDataControlService
}

// NewBetaOrganizationProviderGoogleProjectService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationProviderGoogleProjectService(opts ...option.RequestOption) (r BetaOrganizationProviderGoogleProjectService) {
	r = BetaOrganizationProviderGoogleProjectService{}
	r.Options = opts
	r.DataControls = NewBetaOrganizationProviderGoogleProjectDataControlService(opts...)
	return
}
