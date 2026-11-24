package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaUserProfileProviderService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaUserProfileProviderService] method instead.
type BetaUserProfileProviderService struct {
	Options []option.RequestOption
	Google  BetaUserProfileProviderGoogleService
}

// NewBetaUserProfileProviderService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaUserProfileProviderService(opts ...option.RequestOption) (r BetaUserProfileProviderService) {
	r = BetaUserProfileProviderService{}
	r.Options = opts
	r.Google = NewBetaUserProfileProviderGoogleService(opts...)
	return
}
