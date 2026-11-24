package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaUserProfileProviderGoogleService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaUserProfileProviderGoogleService] method instead.
type BetaUserProfileProviderGoogleService struct {
	Options      []option.RequestOption
	DataControls BetaUserProfileProviderGoogleDataControlService
}

// NewBetaUserProfileProviderGoogleService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaUserProfileProviderGoogleService(opts ...option.RequestOption) (r BetaUserProfileProviderGoogleService) {
	r = BetaUserProfileProviderGoogleService{}
	r.Options = opts
	r.DataControls = NewBetaUserProfileProviderGoogleDataControlService(opts...)
	return
}
