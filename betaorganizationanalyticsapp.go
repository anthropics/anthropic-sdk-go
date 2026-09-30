package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/option"
)

// BetaOrganizationAnalyticsAppService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsAppService] method instead.
type BetaOrganizationAnalyticsAppService struct {
	Options []option.RequestOption
	Chat    BetaOrganizationAnalyticsAppChatService
}

// NewBetaOrganizationAnalyticsAppService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationAnalyticsAppService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsAppService) {
	r = BetaOrganizationAnalyticsAppService{}
	r.Options = opts
	r.Chat = NewBetaOrganizationAnalyticsAppChatService(opts...)
	return
}
