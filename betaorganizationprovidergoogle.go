package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaOrganizationProviderGoogleService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationProviderGoogleService] method instead.
type BetaOrganizationProviderGoogleService struct {
	Options  []option.RequestOption
	Projects BetaOrganizationProviderGoogleProjectService
}

// NewBetaOrganizationProviderGoogleService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationProviderGoogleService(opts ...option.RequestOption) (r BetaOrganizationProviderGoogleService) {
	r = BetaOrganizationProviderGoogleService{}
	r.Options = opts
	r.Projects = NewBetaOrganizationProviderGoogleProjectService(opts...)
	return
}
